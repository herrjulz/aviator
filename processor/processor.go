package processor

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/JulzDiverse/aviator"
	"github.com/JulzDiverse/aviator/filemanager"
	"github.com/JulzDiverse/aviator/printer"
	"github.com/JulzDiverse/aviator/spruce"
	"github.com/pkg/errors"
)

type WriterFunc func([]byte, string) error

type Processor struct {
	spruceClient aviator.SpruceClient
	store        aviator.FileStore
	verbose      bool
	silent       bool
	warnings     []string
}

func NewTestProcessor(spruceClient aviator.SpruceClient, store aviator.FileStore) *Processor {
	return &Processor{
		spruceClient: spruceClient,
		store:        store,
	}
}

func New() *Processor {
	return &Processor{
		store:        filemanager.Store(),
		spruceClient: spruce.New(),
	}
}

func (p *Processor) Process(config []aviator.Spruce) error {
	return p.ProcessWithOpts(config, false, false)
}

func (p *Processor) ProcessVerbose(config []aviator.Spruce) error {
	return p.ProcessWithOpts(config, true, false)
}

func (p *Processor) ProcessSilent(config []aviator.Spruce) error {
	return p.ProcessWithOpts(config, false, true)
}

func (p *Processor) ProcessWithOpts(config []aviator.Spruce, verbose bool, silent bool) error {
	p.verbose, p.silent = verbose, silent
	var err error
	for _, cfg := range config {
		switch mergeType(cfg) {
		case "default":
			err = p.defaultMerge(cfg)
		case "forEach":
			err = p.forEachFileMerge(cfg)
		case "forEachIn":
			err = p.forEachInMerge(cfg)
		case "walkThrough":
			err = p.walk(cfg, "")
		case "walkThroughForAll":
			err = p.forAll(cfg)
		}
	}
	return err
}

func (p *Processor) defaultMerge(cfg aviator.Spruce) error {
	files := p.collectFiles(cfg)
	if err := p.mergeAndWrite(files, cfg, cfg.To); err != nil {
		return errors.Wrap(err, "Spruce Merge FAILED")
	}
	return nil
}

func (p *Processor) forEachFileMerge(cfg aviator.Spruce) error {
	for _, file := range cfg.ForEach.Files {
		mergeFiles := p.collectFiles(cfg)
		fileName, _ := concatFileNameWithPath(file)
		mergeFiles = append(mergeFiles, file)
		targetName := createTargetName(cfg.ToDir, fileName)
		if err := p.mergeAndWrite(mergeFiles, cfg, targetName); err != nil {
			return errors.Wrap(err, "Spruce Merge FAILED")
		}
	}
	return nil
}

func (p *Processor) forEachInMerge(cfg aviator.Spruce) error {
	filePaths, err := ioutil.ReadDir(cfg.ForEach.In)
	if err != nil {
		return err
	}

	regex := getRegexp(cfg.ForEach.Regexp)
	files := p.collectFiles(cfg)
	for _, f := range filePaths {
		if except(cfg.ForEach.Except, f.Name()) {
			p.warnings = append(p.warnings, "SKIPPED: "+f.Name())
			continue
		}
		matched, _ := regexp.MatchString(regex, f.Name())
		if !f.IsDir() && matched {
			prefix := chunk(cfg.ForEach.In)
			mergeFiles := append(files, cfg.ForEach.In+f.Name())
			targetName := createTargetName(cfg.ToDir, fmt.Sprintf("%s_%s", prefix, f.Name()))
			if err := p.mergeAndWrite(mergeFiles, cfg, targetName); err != nil {
				return errors.Wrap(err, "Spruce Merge FAILED")
			}
		} else {
			p.warnings = append(p.warnings, "EXCLUDED BY REGEXP "+regex+": "+cfg.ForEach.In+f.Name())
		}
	}
	return nil
}

func (p *Processor) walk(cfg aviator.Spruce, outer string) error {
	sl := getAllFilesIncludingSubDirs(cfg.ForEach.In)
	regex := getRegexp(cfg.ForEach.Regexp)
	for _, f := range sl {
		filename, parent := concatFileNameWithPath(f)
		match := enableMatching(cfg.ForEach, parent)
		matched, _ := regexp.MatchString(regex, filename)
		if strings.Contains(outer, match) && matched {
			files := p.collectFiles(cfg)
			if outer != "" {
				files = append(files, f, outer)
			} else {
				files = append(files, f)
			}

			if !cfg.ForEach.CopyParents {
				parent = ""
			}

			targetName := createTargetName(cfg.ToDir, filepath.Join(parent, filename))
			if err := p.mergeAndWrite(files, cfg, targetName); err != nil {
				return errors.Wrap(err, "Spruce Merge FAILED")
			}
		}
	}
	return nil
}

func (p *Processor) forAll(cfg aviator.Spruce) error {
	forAll := cfg.ForEach.ForAll
	if forAll != "" {
		files, _ := ioutil.ReadDir(forAll)
		for _, f := range files {
			if !f.IsDir() {
				if err := p.walk(cfg, cfg.ForEach.ForAll+f.Name()); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (p *Processor) mergeAndWrite(files []string, cfg aviator.Spruce, to string) error {
	mergeConf := aviator.MergeConf{
		Files:       files,
		SkipEval:    cfg.SkipEval,
		Prune:       cfg.Prune,
		CherryPicks: cfg.CherryPicks,
	}

	if !p.silent {
		printer.AnsiPrint(mergeConf, to, p.warnings, p.verbose)
	}

	p.warnings = []string{}
	result, err := p.spruceClient.MergeWithOpts(mergeConf)
	if err != nil {
		return errors.Wrap(err, "Spruce Merge FAILED")
	}

	err = p.store.WriteFile(to, result)
	if err != nil {
		return err
	}

	return nil
}

func (p *Processor) collectFiles(cfg aviator.Spruce) []string {
	files := []string{cfg.Base}
	for _, m := range cfg.Merge {
		with := p.collectFilesFromWithSection(m)
		within := p.collectFilesFromWithInSection(m)
		withallin := p.collectFilesFromWithAllInSection(m)
		files = concatStringSlices(files, with, within, withallin)
	}
	return files
}

func (p *Processor) collectFilesFromWithSection(merge aviator.Merge) []string {
	var result []string
	for _, file := range merge.With.Files {
		if merge.With.InDir != "" {
			dir := merge.With.InDir
			file = dir + file
		}

		_, fileExists := p.store.ReadFile(file)
		if !merge.With.Skip || fileExists {
			result = append(result, file)
		} else {
			p.warnings = append(p.warnings, fmt.Sprintf("Skipped non existing file: %s", file))
		}
	}
	return result
}

func (p *Processor) collectFilesFromWithInSection(merge aviator.Merge) []string {
	result := []string{}
	if merge.WithIn != "" {
		within := merge.WithIn
		files, _ := ioutil.ReadDir(within)
		regex := getRegexp(merge.Regexp)
		for _, f := range files {
			if except(merge.Except, f.Name()) {
				continue
			}

			matched, _ := regexp.MatchString(regex, f.Name())
			if !f.IsDir() && matched {
				result = append(result, within+f.Name())
			} else {
				p.warnings = append(p.warnings, "EXCLUDED BY REGEXP "+regex+": "+merge.WithIn+f.Name())
			}
		}
	}
	return result
}

func (p *Processor) collectFilesFromWithAllInSection(merge aviator.Merge) []string {
	result := []string{}
	if merge.WithAllIn != "" {
		allFiles := getAllFilesIncludingSubDirs(merge.WithAllIn)
		regex := getRegexp(merge.Regexp)
		for _, file := range allFiles {
			matched, _ := regexp.MatchString(regex, file)
			if matched {
				result = append(result, file)
			} else {
				p.warnings = append(p.warnings, "EXCLUDED BY REGEXP "+regex+": "+file)
			}
		}
	}
	return result
}
