package processor

import (
	"io/ioutil"
	"regexp"

	"github.com/JulzDiverse/aviator/cockpit"
	"github.com/pkg/errors"
)

//go:generate counterfeiter . SpruceClient
type SpruceClient interface {
	MergeWithOpts(MergeConf) ([]byte, error)
}

type MergeConf struct {
	Files       []string
	Prune       []string
	CherryPicks []string
	SkipEval    bool
}

type Processor struct {
	config       []cockpit.Spruce
	spruceClient SpruceClient
}

func New(spruceClient SpruceClient) *Processor {
	return &Processor{spruceClient: spruceClient}
}

func (p *Processor) Process(config []cockpit.Spruce) ([]byte, error) {
	p.config = config
	for _, cfg := range config {
		var err error
		switch mergeType(cfg) {
		case "default":
			return p.defaultMerge(cfg)
		case "forEach":
		case "forEachIn":
		case "walkThrough":
		case "walkThroughForAll":
		}
		if err != nil {
			return nil, err
		}
	}
	return []byte{}, nil
}

func (p *Processor) defaultMerge(cfg cockpit.Spruce) ([]byte, error) {
	files := p.collectFiles(cfg)
	mergeConf := MergeConf{
		Files:       files,
		Prune:       cfg.Prune,
		SkipEval:    cfg.SkipEval,
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

		if !merge.With.Skip || fileExists(file) { //|| fileExistsInDataStore(file)
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
		regex := getRegexp(merge)
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
		regex := getRegexp(merge)
		for _, file := range allFiles {
			matched, _ := regexp.MatchString(regex, file)
			if matched {
				result = append(result, file)
			}
		}
	}
	return result
}

func (p *Processor) fileExistsInDataStore(file string) {
	//if re.MatchString(path) {
	//matches := re.FindSubmatch([]byte(path))
	//key := string(matches[len(matches)-1])
	//_, ok := spruce.DataStore[key]
	//if ok {
	//return true //return true if dataManager has file
	//}
	//}
}
