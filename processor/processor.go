package processor

import (
	"github.com/JulzDiverse/aviator/spruce"
	"github.com/JulzDiverse/aviator/validator"
)

type SpruceClient interface {
	MergeWithOpts()
}

type SpruceProcessor struct {
	config []validator.Spruce
	spruce SpruceClient
}

func Process(config []validator.Spruce) ([]byte, error) {
	processor := SpruceProcessor{config: config}
	for _, cfg := range config {
		var err error
		switch mergeType(cfg) {
		case "default":
			result, err = processor.defaultMerge(cfg)
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

func mergeType(cfg validator.Spruce) string {
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

func (p *SpruceProcessor) defaultMerge(cfg validator.Spruce) ([]byte, error) {
	files := collectFiles(cfg)
	mergeConf := spruce.MergeOpts{
		Files:       files,
		Prune:       cfg.Prune,
		SkipEval:    cfg.SkipEval,
		CherryPicks: cfg.CherryPicks,
	}
	result, err := p.sprucify(mergeConf, cfg.To)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func collectFiles(cfg validator.Spruce) []string {
	files := []string{cfg.Base}
	for _, val := range cfg.Chain {
		tmp := collectFromMergeSection(val)
		for _, str := range tmp {
			files = append(files, str)
		}
	}
	return files
}

func collectFromMergeSection(merge validator.Merge) []string {
	var result []string
	for _, file := range merge.With.Files {
		if merge.With.InDir != "" {
			dir := merge.With.InDir
			file = dir + file
		}
		if !merge.With.Existing || fileExists(file) {
			result = append(result, file)
		}
	}

	if merge.WithIn != "" {
		within := merge.WithIn
		files, _ := ioutil.ReadDir(within)
		regex := getChainRegexp(merge)
		for _, f := range files {
			if except(merge.Except, f.Name()) {
				continue
			}
			matched, _ := regexp.MatchString(regex, f.Name())
			if matched {
				result = append(result, within+f.Name())
			} else {
				Warnings = append(Warnings, "EXCLUDED BY REGEXP "+regex+": "+merge.WithIn+f.Name())
			}
		}
	}
	return result
}

func (p *SpruceProcessor) sprucify(opts spruce.MergeOpts, fileName string) ([]byte, error) {
	//if !p.Silent {
	//beautifyPrint(opts, fileName)
	//}
	//Warnings = []string{}

	rawYml, err := p.spruce.CmdMergeEval(opts)
	if err != nil {
		return rawYml, err
	}

	resultYml, err := yaml.Marshal(rawYml)
	if err != nil {
		return resultYaml, err
	}

	return nil
}
