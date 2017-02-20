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
	Base      string   `yaml:"base"`
	Files     []string `yaml:"with"`
	FileDir   string   `yaml:"with_in"`
	Prune     []string `yaml:"prune"`
	Folder    string   `yaml:"dir"`
	ForEach   []string `yaml:"for_each"`
	ForEachIn string   `yaml:"for_each_in"`
	Walk      string   `yaml:"walk_through"`
	DestFile  string   `yaml:"to"`
	DestDir   string   `yaml:"to_dir"`
	Regexp    string   `yaml:"regexp"`
}

type FlyConfig struct {
	Config string   `yaml:"config"`
	Vars   []string `yaml:"vars"`
}

func ReadYaml(ymlBytes []byte) Aviator {
	var yml Aviator

	fmt.Printf("%s", ymlBytes)

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
			Walk(conf)
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
		fileName := ConcatFileName(val)
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
			chunked := strings.Split(conf.ForEachIn, "/")
			var prefix string
			if chunked[len(chunked)-1] == "" {
				prefix = chunked[len(chunked)-2]
			} else {
				prefix = chunked[len(chunked)-1]
			}
			cmd = append(cmd, conf.ForEachIn+f.Name())
			SpruceToFile(cmd, conf.DestDir+prefix+"_"+f.Name())
		}
	}
}

func Walk(conf SpruceConfig) {
	sl := []string{}
	err := filepath.Walk(conf.Walk, fillSliceWithFiles(&sl))
	if err != nil {
		log.Fatal(err)
	}
	regex := getRegexp(conf)

	for _, f := range sl {
		fileName := ConcatFileName(f)
		matched, _ := regexp.MatchString(regex, fileName)
		if matched {
			cmd := CreateSpruceCommand(conf)
			cmd = append(cmd, f)
			SpruceToFile(cmd, conf.DestDir+fileName)
		}
	}
}

func getRegexp(conf SpruceConfig) string {
	regex := ".*"
	if conf.Regexp != "" {
		regex = conf.Regexp
	}
	return regex
}

func ConcatFileName(path string) string {
	chunked := strings.Split(path, "/")
	fileName := chunked[len(chunked)-2] + "_" + chunked[len(chunked)-1]
	return fileName
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
	for _, file := range spruce.Files {
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
	cmd.Wait()
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
