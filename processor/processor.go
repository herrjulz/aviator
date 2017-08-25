package processor

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/JulzDiverse/aviator/cockpit"
	"github.com/JulzDiverse/aviator/filemanager"
	"github.com/pkg/errors"
)

//go:generate counterfeiter . SpruceClient
type SpruceClient interface {
	MergeWithOpts(MergeConf) ([]byte, error)
}

//go:generate counterfeiter . FileStore
type FileStore interface {
	ReadFile(string) ([]byte, bool)
	WriteFile(string, []byte) error
}

type MergeConf struct {
	Files       []string
	Prune       []string
	CherryPicks []string
	SkipEval    bool
	Warnings    []string
	To          string
}

type WriterFunc func([]byte, string) error

type Processor struct {
	spruceClient SpruceClient
	store        FileStore
	verbose      bool
	silent       bool
}

func NewTestProcessor(spruceClient SpruceClient, store FileStore) *Processor {
	return &Processor{
		spruceClient: spruceClient,
		store:        store,
	}
}

type ProcessorFactory func(*Processor, bool, bool)

func Create(fn ProcessorFactory, silent bool, verbose bool) *Processor {
	var p Processor
	fn(&p, silent, verbose)
	return &p
}

func AviatorDefault(p *Processor, s bool, v bool) {
	p.store = filemanager.Store()
	p.verbose = v
	p.silent = s
}

func (p *Processor) Process(config []cockpit.Spruce) error {
	for _, cfg := range config {
		switch mergeType(cfg) {
		case "default":
			return p.defaultMerge(cfg)
		case "forEach":
			return p.forEachFileMerge(cfg)
		case "forEachIn":
			return p.forEachInMerge(cfg)
		case "walkThrough":
			return p.walk(cfg, "")
		case "walkThroughForAll":
			return p.forAll(cfg)
		}
	}
	return nil
}

func (p *Processor) defaultMerge(cfg cockpit.Spruce) error {
	files := p.collectFiles(cfg)
	if err := p.mergeAndWrite(files, cfg, cfg.To); err != nil {
		return errors.Wrap(err, "Spruce Merge FAILED")
	}
	return nil
}

func (p *Processor) forEachFileMerge(cfg cockpit.Spruce) error {
	for _, file := range cfg.ForEach.Files {
		mergeFiles := p.collectFiles(cfg)
		fileName, _ := concatFileNameWithPath(file)
		mergeFiles = append(mergeFiles, file)
		if err := p.mergeAndWrite(mergeFiles, cfg, filepath.Join(cfg.ToDir, fileName)); err != nil {
			return errors.Wrap(err, "Spruce Merge FAILED")
		}
	}
	return nil
}

func (p *Processor) forEachInMerge(cfg cockpit.Spruce) error {
	filePaths, err := ioutil.ReadDir(cfg.ForEach.In)
	if err != nil {
		return err
	}

	regex := getRegexp(cfg.ForEach.Regexp)
	files := p.collectFiles(cfg)
	for _, f := range filePaths {
		if except(cfg.ForEach.Except, f.Name()) {
			//Warnings = append(Warnings, "SKIPPED: "+f.Name())
			continue
		}
		matched, _ := regexp.MatchString(regex, f.Name())
		if !f.IsDir() && matched {
			prefix := chunk(cfg.ForEach.In)
			mergeFiles := append(files, cfg.ForEach.In+f.Name())
			if err := p.mergeAndWrite(mergeFiles, cfg, filepath.Join(cfg.ToDir, fmt.Sprintf("%s_%s", prefix, f.Name()))); err != nil {
				return errors.Wrap(err, "Spruce Merge FAILED")
			}
		} else {
			//Warnings = append(Warnings, "EXCLUDED BY REGEXP "+regex+": "+conf.ForEachIn+f.Name())
		}
	}
	return nil
}

func (p *Processor) walk(cfg cockpit.Spruce, outer string) error {
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

			if err := p.mergeAndWrite(files, cfg, filepath.Join(cfg.ToDir, parent, filename)); err != nil {
				return errors.Wrap(err, "Spruce Merge FAILED")
			}
		}
	}
	return nil
}

func (p *Processor) forAll(cfg cockpit.Spruce) error {
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

func (p *Processor) mergeAndWrite(files []string, cfg cockpit.Spruce, to string) error {
	mergeConf := MergeConf{
		Files:       files,
		SkipEval:    cfg.SkipEval,
		Prune:       cfg.Prune,
		CherryPicks: cfg.CherryPicks,
	}
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

func (p *Processor) collectFiles(cfg cockpit.Spruce) []string {
	files := []string{cfg.Base}
	for _, m := range cfg.Merge {
		with := p.collectFilesFromWithSection(m)
		within := p.collectFilesFromWithInSection(m)
		withallin := p.collectFilesFromWithAllInSection(m)
		files = concatStringSlices(files, with, within, withallin)
	}
	return files
}

func (p *Processor) collectFilesFromWithSection(merge cockpit.Merge) []string {
	var result []string
	for _, file := range merge.With.Files {
		if merge.With.InDir != "" {
			dir := merge.With.InDir
			file = dir + file
		}

		_, fileExists := p.store.ReadFile(file)
		if !merge.With.Skip || fileExists {
			result = append(result, file)
		}
	}
	return result
}

func (p *Processor) collectFilesFromWithInSection(merge cockpit.Merge) []string {
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
			}
			//else {
			//Warnings = append(Warnings, "EXCLUDED BY REGEXP "+regex+": "+merge.WithIn+f.Name())
			//}
		}
	}
	return result
}

func (p *Processor) collectFilesFromWithAllInSection(merge cockpit.Merge) []string {
	result := []string{}
	if merge.WithAllIn != "" {
		allFiles := getAllFilesIncludingSubDirs(merge.WithAllIn)
		regex := getRegexp(merge.Regexp)
		for _, file := range allFiles {
			matched, _ := regexp.MatchString(regex, file)
			if matched {
				result = append(result, file)
			}
		}
	}
	return result
}
