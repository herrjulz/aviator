package processor

import (
	"io/ioutil"
	"os"

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

func mergeType(cfg cockpit.Spruce) string {
	if cfg.ForEachIn == "" && len(cfg.ForEach) == 0 && cfg.WalkThrough == "" {
		return "default"
	}
	if len(cfg.ForEach) != 0 {
		return "forEach"
	}
	if cfg.ForEachIn != "" {
		return "forEachIn"
	}
	if cfg.WalkThrough != "" {
		if cfg.ForAll != "" {
			return "walkThrough"
		} else {
			return "walkThroughForAll"
		}
	}
	return ""
}

func (p *Processor) defaultMerge(cfg cockpit.Spruce) ([]byte, error) {
	files := collectFiles(cfg)
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

func collectFiles(cfg cockpit.Spruce) []string {
	files := []string{cfg.Base}
	for _, m := range cfg.Merge {
		with := collectFilesFromWithSection(m)
		within := collectFilesFromWithInSection(m)
		files = concatStringSlices(files, with, within)
	}
	return files
}

func collectFilesFromWithSection(merge cockpit.Merge) []string {
	var result []string
	for _, file := range merge.With.Files {
		if merge.With.InDir != "" {
			dir := merge.With.InDir
			file = dir + file
		}

		if !merge.With.Existing || fileExists(file) { //|| fileExistsInDataStore(file)
			result = append(result, file)
		}
	}
	return result
}

func collectFilesFromWithInSection(merge cockpit.Merge) []string {
	result := []string{}
	if merge.WithIn != "" {
		within := merge.WithIn
		files, _ := ioutil.ReadDir(within)
		//regex := regexp(merge)
		for _, f := range files {
			if except(merge.Except, f.Name()) {
				continue
			}
			//matched, _ := regexp.MatchString(regex, f.Name())
			//if matched {
			//result = append(result, within+f.Name())
			//} else {
			//Warnings = append(Warnings, "EXCLUDED BY REGEXP "+regex+": "+merge.WithIn+f.Name())
			//}
			if !f.IsDir() {
				result = append(result, within+f.Name())
			}
		}
	}
	return result
}

func fileExistsInDataStore(file string) {
	//if re.MatchString(path) {
	//matches := re.FindSubmatch([]byte(path))
	//key := string(matches[len(matches)-1])
	//_, ok := spruce.DataStore[key]
	//if ok {
	//return true //return true if dataManager has file
	//}
	//}
}

func except(except []string, file string) bool {
	for _, f := range except {
		if f == file {
			return true
		}
	}
	return false
}

func regexp(merge cockpit.Merge) string {
	regex := ".*"
	if merge.Regexp != "" {
		regex = merge.Regexp
	}
	return regex
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func concatStringSlices(sl1 []string, sls ...[]string) []string {
	for _, sl := range sls {
		for _, s := range sl {
			sl1 = append(sl1, s)
		}
	}
	return sl1
}

//}

//func (p *SpruceProcessor) sprucify(opts spruce.MergeOpts, fileName string) ([]byte, error) {
////if !p.Silent {
////beautifyPrint(opts, fileName)
////}
////Warnings = []string{}

//rawYml, err := p.spruce.CmdMergeEval(opts)
//if err != nil {
//return rawYml, err
//}

//resultYml, err := yaml.Marshal(rawYml)
//if err != nil {
//return resultYaml, err
//}

//return nil
//}
