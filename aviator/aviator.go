package aviator

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

type Aviator struct {
	Spruce []SpruceConfig `yaml:"spruce"`
	Fly    FlyConfig      `yaml:"fly"`
}

type SpruceConfig struct {
	Base           string   `yaml:"base"`
	Prune          []string `yaml:"prune"`
	Chain          []Chain  `yaml:"chain"`
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
	Files []string `yaml:"files"`
	InDir string   `yaml:"in_dir"`
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
			straight(conf)
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

func straight(conf SpruceConfig) {
	cmd := ProcessChain(conf)
	dest := resolveVar(conf.DestFile)
	SpruceToFile(cmd, dest)
}

func ProcessChain(conf SpruceConfig) []string {
	cmd := createBaseCommand(conf)
	for _, val := range conf.Chain {
		tmp := CreateSpruceCommand(val)
		for _, str := range tmp {
			cmd = append(cmd, str)
		}
	}
	return cmd
}

func ForEachFile(conf SpruceConfig) {
	dest := resolveVar(conf.DestFile)

	for _, val := range conf.ForEach {
		// cmd := CreateSpruceCommand(conf.Chain)
		val = resolveVar(val)
		cmd := ProcessChain(conf)
		fileName, _ := ConcatFileName(val)
		cmd = append(cmd, val)
		SpruceToFile(cmd, dest+fileName)
	}
}

func ForEachIn(conf SpruceConfig) {
	forEachIn := resolveVar(conf.ForEachIn)
	dest := resolveVar(conf.DestFile)
	files, _ := ioutil.ReadDir(forEachIn)
	regex := getRegexp(conf)
	for _, f := range files {
		// cmd := CreateSpruceCommand(conf)
		cmd := ProcessChain(conf)
		matched, _ := regexp.MatchString(regex, f.Name())
		if matched {
			prefix := Chunk(conf.ForEachIn)
			cmd = append(cmd, forEachIn+f.Name())
			SpruceToFile(cmd, dest+prefix+"_"+f.Name())
		}
	}
}

func ForEachInner(conf SpruceConfig, outer string) {
	// forEachIn := resolveVar(conf.ForEachIn)
	dest := resolveVar(conf.DestFile)
	files, _ := ioutil.ReadDir(conf.ForEachIn)
	regex := getRegexp(conf)
	for _, f := range files {
		// cmd := CreateSpruceCommand(conf)
		cmd := ProcessChain(conf)
		matched, _ := regexp.MatchString(regex, f.Name())
		if matched {
			prefix := Chunk(conf.ForEachIn)
			cmd = append(cmd, conf.ForEachIn+f.Name())
			cmd = append(cmd, outer)
			SpruceToFile(cmd, dest+prefix+"_"+f.Name())
		}
	}
}

func ForAll(conf SpruceConfig) {
	forAll := resolveVar(conf.ForAll)
	if forAll != "" {
		files, _ := ioutil.ReadDir(forAll)
		for _, f := range files {
			//ForEachInner(conf, conf.ForAll+f.Name())
			Walk(conf, forAll+f.Name())
		}
	}
}

func Walk(conf SpruceConfig, outer string) {
	dest := resolveVar(conf.DestFile)
	walk := resolveVar(conf.Walk)
	sl := getAllFilesInSubDirs(walk)
	regex := getRegexp(conf)

	for _, f := range sl {
		filename, parent := ConcatFileName(f)
		match := isMatchingEnabled(conf, parent)
		if strings.Contains(outer, match) {
			matched, _ := regexp.MatchString(regex, filename)
			if matched {
				// cmd := CreateSpruceCommand(conf)
				cmd := ProcessChain(conf)
				cmd = append(cmd, f)
				cmd = append(cmd, outer)
				if conf.CopyParents {
					CreateDir(dest + parent)
				} else {
					parent = ""
				}
				SpruceToFile(cmd, dest+parent+"/"+filename)
			}
		}
	}
}

func createBaseCommand(conf SpruceConfig) []string {
	spruceCmd := []string{"--concourse", "merge"}
	for _, prune := range conf.Prune {
		spruceCmd = append(spruceCmd, "--prune", prune)
	}
	base := resolveVar(conf.Base)
	spruceCmd = append(spruceCmd, base)
	return spruceCmd
}

func CreateSpruceCommand(chain Chain) []string {
	var spruceCmd []string
	for _, file := range chain.With.Files {
		file = resolveVar(file)
		if chain.With.InDir != "" {
			dir := resolveVar(chain.With.InDir)
			file = dir + file
		}
		spruceCmd = append(spruceCmd, file)
	}

	if chain.WithIn != "" {
		within := resolveVar(chain.WithIn)
		files, _ := ioutil.ReadDir(within)
		regex := getChainRegexp(chain)
		for _, f := range files {
			matched, _ := regexp.MatchString(regex, f.Name())
			if matched {
				spruceCmd = append(spruceCmd, within+f.Name())
			}
		}
	}

	return spruceCmd
}

func SpruceToFile(argv []string, fileName string) {
	cmd := exec.Command("spruce", argv...)
	fmt.Println("EXEC SPRUCE:", cmd.Args, "to", fileName)
	outfile, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer outfile.Close()

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	printStderr(cmd)

	writer := bufio.NewWriter(outfile)
	defer writer.Flush()

	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	go io.Copy(writer, stdoutPipe)

	err = cmd.Wait()
	if err != nil {
		os.Exit(1)
	}
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

// func Walk(conf SpruceConfig) {
// 	sl := getAllFilesInSubDirs(conf.Walk)
// 	regex := getRegexp(conf)
//
// 	for _, f := range sl {
// 		filename, parent := ConcatFileName(f)
// 		matched, _ := regexp.MatchString(regex, filename)
// 		if matched {
// 			cmd := CreateSpruceCommand(conf)
// 			cmd = append(cmd, f)
// 			CreateDir(conf.DestDir + parent)
// 			SpruceToFile(cmd, conf.DestDir+parent+"/"+filename)
// 		}
// 	}
// }
