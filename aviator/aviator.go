package aviator

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/JulzDiverse/aviator/spruce"

	"gopkg.in/yaml.v2"
)

type Aviator struct {
	Spruce []SpruceConfig `yaml:"spruce"`
	Fly    FlyConfig      `yaml:"fly"`
}

type SpruceConfig struct {
	Base           string   `yaml:"base"`
	Prune          []string `yaml:"prune"`
	Chain          []Chain  `yaml:"merge"`
	WithIn         string   `yaml:"with_in"`
	Folder         string   `yaml:"dir"`
	ForEach        []string `yaml:"for_each"`
	ForEachIn      string   `yaml:"for_each_in"`
	Walk           string   `yaml:"walk_through"`
	ForAll         string   `yaml:"for_all"`
	CopyParents    bool     `yaml:"copy_parents"`
	EnableMatching bool     `yaml:"enable_matching"`
	DestFile       string   `yaml:"to"`
	DestDir        string   `yaml:"to_dir"`
	Regexp         string   `yaml:"regexp"`
}

type Chain struct {
	With   With   `yaml:"with"`
	WithIn string `yaml:"with_in"`
	Regexp string `yaml:"regexp"`
}

type With struct {
	Files    []string `yaml:"files"`
	InDir    string   `yaml:"in_dir"`
	Existing bool     `yaml:"skip_non_existing"`
}

type FlyConfig struct {
	Config string   `yaml:"config"`
	Vars   []string `yaml:"vars"`
}

func ReadYaml(ymlBytes []byte) Aviator {
	var yml Aviator

	// fmt.Printf("%s", ymlBytes)

	err := yaml.Unmarshal(ymlBytes, &yml)
	if err != nil {
		panic(err)
	}

	return yml
}

func FlyPipeline(fly FlyConfig, target string, pipeline string) {

	flyCmd := []string{"-t", target, "set-pipeline", "-p", pipeline, "-c", fly.Config}
	for _, val := range fly.Vars {
		flyCmd = append(flyCmd, "-l", val)
	}

	cmd := exec.Command("fly", flyCmd...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failed to run fly. %s\n", err.Error())
		os.Exit(1)
	}
}

func ProcessSprucePlan(spruce []SpruceConfig) {
	for _, conf := range spruce {

		verifySpruceConfig(conf)

		if conf.ForEachIn == "" && len(conf.ForEach) == 0 && conf.Walk == "" {
			simpleMerge(conf)
		}
		if len(conf.ForEach) != 0 {
			ForEachFile(conf)
		}
		if conf.ForEachIn != "" {
			ForEachIn(conf)
		}
		if conf.Walk != "" {
			if conf.ForAll != "" {
				ForAll(conf)
			} else {
				Walk(conf, "")
			}
		}
	}
}

func simpleMerge(conf SpruceConfig) {
	files := collectFiles(conf)
	mergeConf := spruce.MergeOpts{
		Files: files,
		Prune: conf.Prune,
	}
	spruceToFile(mergeConf, conf.DestFile)
}

func collectFiles(conf SpruceConfig) []string {
	files := []string{}
	for _, val := range conf.Chain {
		tmp := collectFromMergeSection(val)
		for _, str := range tmp {
			files = append(files, str)
		}
	}
	return files
}

func ForEachFile(conf SpruceConfig) {
	for _, val := range conf.ForEach {
		files := collectFiles(conf)
		fileName, _ := ConcatFileName(val)
		files = append(files, val)
		mergeConf := spruce.MergeOpts{
			Files: files,
			Prune: conf.Prune,
		}
		spruceToFile(mergeConf, conf.DestDir+fileName)
	}
}

func ForEachIn(conf SpruceConfig) {
	filePaths, _ := ioutil.ReadDir(conf.ForEachIn)
	regex := getRegexp(conf)
	for _, f := range filePaths {
		files := collectFiles(conf)
		matched, _ := regexp.MatchString(regex, f.Name())
		if matched {
			prefix := Chunk(conf.ForEachIn)
			files = append(files, conf.ForEachIn+f.Name())
			mergeConf := spruce.MergeOpts{
				Files: files,
				Prune: conf.Prune,
			}
			spruceToFile(mergeConf, conf.DestDir+prefix+"_"+f.Name())
		}
	}
}

func ForEachInner(conf SpruceConfig, outer string) {
	filePaths, _ := ioutil.ReadDir(conf.ForEachIn)
	regex := getRegexp(conf)
	for _, f := range filePaths {
		files := collectFiles(conf)
		matched, _ := regexp.MatchString(regex, f.Name())
		if matched {
			prefix := Chunk(conf.ForEachIn)
			files = append(files, conf.ForEachIn+f.Name())
			files = append(files, outer)
			mergeConf := spruce.MergeOpts{
				Files: files,
				Prune: conf.Prune,
			}
			spruceToFile(mergeConf, conf.DestDir+prefix+"_"+f.Name())
		}
	}
}

func ForAll(conf SpruceConfig) {
	if conf.ForAll != "" {
		files, _ := ioutil.ReadDir(conf.ForAll)
		for _, f := range files {
			Walk(conf, conf.ForAll+f.Name())
		}
	}
}

func Walk(conf SpruceConfig, outer string) {
	sl := getAllFilesInSubDirs(conf.Walk)
	regex := getRegexp(conf)

	for _, f := range sl {
		filename, parent := ConcatFileName(f)
		match := isMatchingEnabled(conf, parent)
		if strings.Contains(outer, match) {
			matched, _ := regexp.MatchString(regex, filename)
			if matched {
				files := collectFiles(conf)
				files = append(files, f)
				files = append(files, outer)
				if conf.CopyParents {
					CreateDir(conf.DestDir + parent)
				} else {
					parent = ""
				}
				mergeConf := spruce.MergeOpts{
					Files: files,
					Prune: conf.Prune,
				}
				spruceToFile(mergeConf, conf.DestDir+parent+"/"+filename)
			}
		}
	}
}

func collectFromMergeSection(chain Chain) []string {
	var result []string
	for _, file := range chain.With.Files {
		if chain.With.InDir != "" {
			dir := chain.With.InDir
			file = dir + file
		}
		if !chain.With.Existing || fileExists(file) {
			result = append(result, file)
		}
	}

	if chain.WithIn != "" {
		within := chain.WithIn
		files, _ := ioutil.ReadDir(within)
		regex := getChainRegexp(chain)
		for _, f := range files {
			matched, _ := regexp.MatchString(regex, f.Name())
			if matched {
				result = append(result, within+f.Name())
			}
		}
	}

	return result
}

func spruceToFile(opts spruce.MergeOpts, fileName string) {
	beautifyPrint(opts, fileName)
	rawYml, _ := spruce.CmdMergeEval(opts)
	resultYml, _ := yaml.Marshal(rawYml)
	ioutil.WriteFile(fileName, resultYml, 0644)
}

func Cleanup(path string) {
	d, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		panic(err)
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(path, name))
		if err != nil {
			panic(err)
		}
	}
}
