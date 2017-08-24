package processor

import (
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/JulzDiverse/aviator/cockpit"
	"github.com/pkg/errors"
)

//go:generate counterfeiter . SpruceClient
type SpruceClient interface {
	MergeWithOpts(MergeConf) ([]byte, error)
}

//go:generate counterfeiter . FileStore
type FileStore interface {
	GetFile(string) ([]byte, bool)
	SetFile([]byte, string)
}

type MergeConf struct {
	Files       []string
	Prune       []string
	CherryPicks []string
	SkipEval    bool
	Warnings    []string
}

type Processor struct {
	config       []cockpit.Spruce
	spruceClient SpruceClient
	mergeOpts    MergeConf
	store        FileStore
}

func New(spruceClient SpruceClient, store FileStore) *Processor {
	return &Processor{
		spruceClient: spruceClient,
		store:        store,
	}
}

func (p *Processor) Process(config []cockpit.Spruce) ([][]byte, error) {
	for _, cfg := range config {
		var err error
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
		if err != nil {
			return nil, err
		}
	}
	return [][]byte{}, nil
}

func (p *Processor) defaultMerge(cfg cockpit.Spruce) ([][]byte, error) {
	files := p.collectFiles(cfg)
	result, err := p.merge(files, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "Spruce Merge FAILED")
	}
	return [][]byte{result}, nil
}

func (p *Processor) forEachFileMerge(cfg cockpit.Spruce) ([][]byte, error) {
	mergedFiles := [][]byte{}
	for _, file := range cfg.ForEach.Files {
		mergeFiles := p.collectFiles(cfg)
		//fileName, _ := concatFileNameWithPath(file) --> will be part of cmd
		mergeFiles = append(mergeFiles, file)
		result, err := p.merge(mergeFiles, cfg)
		if err != nil {
			return nil, errors.Wrap(err, "Spruce Merge FAILED")
		}
		mergedFiles = append(mergedFiles, result)
	}
	return mergedFiles, nil
}

func (p *Processor) forEachInMerge(cfg cockpit.Spruce) ([][]byte, error) {
	mergedFiles := [][]byte{}
	filePaths, err := ioutil.ReadDir(cfg.ForEach.In)
	if err != nil {
		return nil, err
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
			//prefix := chunk(cfg.ForEach.In)
			mergeFiles := append(files, cfg.ForEach.In+f.Name())
			result, err := p.merge(mergeFiles, cfg)
			if err != nil {
				return nil, errors.Wrap(err, "Spruce Merge FAILED")
			}

			mergedFiles = append(mergedFiles, result)
		} else {
			//Warnings = append(Warnings, "EXCLUDED BY REGEXP "+regex+": "+conf.ForEachIn+f.Name())
		}
	}
	return mergedFiles, nil
}

func (p *Processor) walk(cfg cockpit.Spruce, outer string) ([][]byte, error) {
	mergedFiles := [][]byte{}
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
			result, err := p.merge(files, cfg)
			if err != nil {
				return nil, errors.Wrap(err, "Spruce Merge FAILED")
			}
			mergedFiles = append(mergedFiles, result)
		}
	}
	return mergedFiles, nil
}

func (p *Processor) forAll(cfg cockpit.Spruce) ([][]byte, error) {
	mergedFiles := [][]byte{}
	forAll := cfg.ForEach.ForAll
	if forAll != "" {
		files, _ := ioutil.ReadDir(forAll)
		for _, f := range files {
			if !f.IsDir() {
				results, err := p.walk(cfg, cfg.ForEach.ForAll+f.Name())
				if err != nil {
					return nil, err
				}
				mergedFiles = concatResults(mergedFiles, results)
			}
		}
	}
	return mergedFiles, nil
}

func (p *Processor) merge(files []string, cfg cockpit.Spruce) ([]byte, error) {
	mergeConf := MergeConf{
		Files:       files,
		SkipEval:    cfg.SkipEval,
		Prune:       cfg.Prune,
		CherryPicks: cfg.CherryPicks,
	}
	result, err := p.spruceClient.MergeWithOpts(mergeConf)
	if err != nil {
		return nil, errors.Wrap(err, "Spruce Merge FAILED")
	}
	return result, nil
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

		_, storeHasFile := p.store.GetFile(file)
		if !merge.With.Skip || fileExists(file) || storeHasFile {
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
