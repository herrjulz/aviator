package aviator

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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
	With           With     `yaml:"with"`
	FileDir        string   `yaml:"with_in"`
	Prune          []string `yaml:"prune"`
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

func ProcessSpruceChain(spruce []SpruceConfig) {
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

func verifySpruceConfig(conf SpruceConfig) {
	if (len(conf.ForEach) != 0 && conf.DestFile != "") ||
		(conf.ForEachIn != "" && conf.DestFile != "") ||
		(conf.Walk != "" && conf.DestFile != "") {
		fmt.Println("'for_each', 'for_each_in', and 'walk_through' in combination with 'to_file' is not allowed: Cannot spruce multiple YAMLs to one destiantion file. ")
		os.Exit(1)
	}
	if len(conf.ForEach) != 0 && conf.ForEachIn != "" {
		fmt.Println("'for_each' in combination with 'for_each_in' is not allowed: Either you want to spruce merge with specific files or files within a directiory. ")
		os.Exit(1)
	}
}

func straight(conf SpruceConfig) {
	cmd := CreateSpruceCommand(conf)
	SpruceToFile(cmd, conf.DestFile)
}

func ForEachFile(conf SpruceConfig) {
	for _, val := range conf.ForEach {
		cmd := CreateSpruceCommand(conf)
		fileName, _ := ConcatFileName(val)
		cmd = append(cmd, val)
		SpruceToFile(cmd, conf.DestDir+fileName)
	}
}

func ForEachIn(conf SpruceConfig) {
	files, _ := ioutil.ReadDir(conf.ForEachIn)
	regex := getRegexp(conf)
	for _, f := range files {
		cmd := CreateSpruceCommand(conf)
		matched, _ := regexp.MatchString(regex, f.Name())
		if matched {
			prefix := Chunk(conf.ForEachIn)
			cmd = append(cmd, conf.ForEachIn+f.Name())
			SpruceToFile(cmd, conf.DestDir+prefix+"_"+f.Name())
		}
	}
}

func ForEachInner(conf SpruceConfig, outer string) {
	files, _ := ioutil.ReadDir(conf.ForEachIn)
	regex := getRegexp(conf)
	for _, f := range files {
		cmd := CreateSpruceCommand(conf)
		matched, _ := regexp.MatchString(regex, f.Name())
		if matched {
			prefix := Chunk(conf.ForEachIn)
			cmd = append(cmd, conf.ForEachIn+f.Name())
			cmd = append(cmd, outer)
			SpruceToFile(cmd, conf.DestDir+prefix+"_"+f.Name())
		}
	}
}

func ForAll(conf SpruceConfig) {
	if conf.ForAll != "" {
		files, _ := ioutil.ReadDir(conf.ForAll)
		for _, f := range files {
			//ForEachInner(conf, conf.ForAll+f.Name())
			Walk(conf, conf.ForAll+f.Name())
		}
	}
}

func isMatchingEnabled(conf SpruceConfig, match string) string {
	if !conf.EnableMatching {
		match = ""
	}
	return match
}

func Chunk(path string) string {
	chunked := strings.Split(path, "/")
	var prefix string
	if chunked[len(chunked)-1] == "" {
		prefix = chunked[len(chunked)-2]
	} else {
		prefix = chunked[len(chunked)-1]
	}
	return prefix
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
				cmd := CreateSpruceCommand(conf)
				cmd = append(cmd, f)
				cmd = append(cmd, outer)
				if conf.CopyParents {
					CreateDir(conf.DestDir + parent)
				} else {
					parent = ""
				}
				SpruceToFile(cmd, conf.DestDir+parent+"/"+filename)
			}
		}
	}
}

func CreateDir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0711)
	}
}

func getAllFilesInSubDirs(path string) []string {
	sl := []string{}
	err := filepath.Walk(path, fillSliceWithFiles(&sl))
	if err != nil {
		log.Fatal(err)
	}
	return sl
}

func getRegexp(conf SpruceConfig) string {
	regex := ".*"
	if conf.Regexp != "" {
		regex = conf.Regexp
	}
	return regex
}

func ConcatFileName(path string) (string, string) {
	chunked := strings.Split(path, "/")
	fileName := chunked[len(chunked)-2] + "_" + chunked[len(chunked)-1]
	parent := chunked[len(chunked)-2]
	return fileName, parent
}

func fillSliceWithFiles(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			*files = append(*files, path)
		}
		return nil
	}
}

func CreateSpruceCommand(spruce SpruceConfig) []string {
	spruceCmd := []string{"--concourse", "merge"}
	for _, prune := range spruce.Prune {
		spruceCmd = append(spruceCmd, "--prune", prune)
	}
	spruceCmd = append(spruceCmd, spruce.Base)
	for _, file := range spruce.With.Files {
		if spruce.With.InDir != "" {
			file = spruce.With.InDir + file
		}
		spruceCmd = append(spruceCmd, file)
	}

	if spruce.FileDir != "" {
		files, _ := ioutil.ReadDir(spruce.FileDir)
		regex := getRegexp(spruce)
		for _, f := range files {
			matched, _ := regexp.MatchString(regex, f.Name())
			if matched {
				spruceCmd = append(spruceCmd, spruce.FileDir+f.Name())
			}
		}
	}

	return spruceCmd
}

func SpruceToFile(argv []string, fileName string) {
	cmd := exec.Command("spruce", argv...)
	fmt.Println("EXEC SPRUCE:", cmd.Args)
	outfile, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer outfile.Close()

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

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

func scanStderr(cmd exec.Cmd) {
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}

	errorScanner := bufio.NewScanner(stderrPipe)
	go func() {
		for errorScanner.Scan() {
			fmt.Printf("%s\n", errorScanner.Text())
		}
	}()
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
