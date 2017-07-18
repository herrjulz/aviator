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
	With   With     `yaml:"with"`
	WithIn string   `yaml:"with_in"`
	Except []string `yaml:"except"`
	Regexp string   `yaml:"regexp"`
}

type With struct {
	Files    []string `yaml:"files"`
	InDir    string   `yaml:"in_dir"`
	Existing bool     `yaml:"skip_non_existing"`
}

type FlyConfig struct {
	Name   string   `yaml:"name"`
	Target string   `yaml:"target"`
	Config string   `yaml:"config"`
	Vars   []string `yaml:"vars"`
}

func ReadYaml(ymlBytes []byte) Aviator {
	var yml Aviator

	ymlBytes = quoteBraces(ymlBytes)
	err := yaml.Unmarshal(ymlBytes, &yml)
	if err != nil {
		panic(err)
	}

	return yml
}

var quoteRegex = `\{\{([-\w\p{L}]+)\}\}`
var re = regexp.MustCompile("(" + quoteRegex + ")")

func quoteBraces(input []byte) []byte {
	return re.ReplaceAll(input, []byte("\"$1\""))
}

func FlyPipeline(fly FlyConfig) {

	flyCmd := []string{"-t", fly.Target, "set-pipeline", "-p", fly.Name, "-c", fly.Config}
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

func ProcessSprucePlan(spruce []SpruceConfig) error {
	for _, conf := range spruce {

		verifySpruceConfig(conf)

		if conf.ForEachIn == "" && len(conf.ForEach) == 0 && conf.Walk == "" {
			err := simpleMerge(conf)
			if err != nil {
				return err
			}
		}
		if len(conf.ForEach) != 0 {
			err := ForEachFile(conf)
			if err != nil {
				return err
			}
		}
		if conf.ForEachIn != "" {
			err := ForEachIn(conf)
			if err != nil {
				return err
			}
		}
		if conf.Walk != "" {
			if conf.ForAll != "" {
				err := ForAll(conf)
				if err != nil {
					return err
				}
			} else {
				err := Walk(conf, "")
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func simpleMerge(conf SpruceConfig) error {
	files := collectFiles(conf)
	mergeConf := spruce.MergeOpts{
		Files: files,
		Prune: conf.Prune,
	}
	err := spruceToFile(mergeConf, conf.DestFile)
	if err != nil {
		return err
	}
	return nil
}

func collectFiles(conf SpruceConfig) []string {
	files := []string{conf.Base}
	for _, val := range conf.Chain {
		tmp := collectFromMergeSection(val)
		for _, str := range tmp {
			files = append(files, str)
		}
	}
	return files
}

func ForEachFile(conf SpruceConfig) error {
	for _, val := range conf.ForEach {
		files := collectFiles(conf)
		fileName, _ := ConcatFileName(val)
		files = append(files, val)
		mergeConf := spruce.MergeOpts{
			Files: files,
			Prune: conf.Prune,
		}
		err := spruceToFile(mergeConf, conf.DestDir+fileName)
		if err != nil {
			return err
		}

	}
	return nil
}

func ForEachIn(conf SpruceConfig) error {
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
			err := spruceToFile(mergeConf, conf.DestDir+prefix+"_"+f.Name())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func ForEachInner(conf SpruceConfig, outer string) error {
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
			err := spruceToFile(mergeConf, conf.DestDir+prefix+"_"+f.Name())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func ForAll(conf SpruceConfig) error {
	if conf.ForAll != "" {
		files, _ := ioutil.ReadDir(conf.ForAll)
		for _, f := range files {
			err := Walk(conf, conf.ForAll+f.Name())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func Walk(conf SpruceConfig, outer string) error {
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
				err := spruceToFile(mergeConf, conf.DestDir+parent+"/"+filename)
				if err != nil {
					return err
				}

			}
		}
	}
	return nil
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
			if except(chain.Except, f.Name()) {
				continue
			}
			matched, _ := regexp.MatchString(regex, f.Name())
			if matched {
				result = append(result, within+f.Name())
			}
		}
	}

	return result
}

func except(except []string, file string) bool {
	for _, f := range except {
		if f == file {
			return true
		}
	}
	return false
}

func spruceToFile(opts spruce.MergeOpts, fileName string) error {
	beautifyPrint(opts, fileName)
	rawYml, err := spruce.CmdMergeEval(opts)
	if err != nil {
		return err
	}

	resultYml, err := yaml.Marshal(rawYml)
	if err != nil {
		return err
	}

	spruce.WriteYamlToPathOrStore(fileName, resultYml)
	return nil
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
