package aviator

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

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

func CreateDir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0711)
	}
}

// helper function for Walk()
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

func getChainRegexp(conf Chain) string {
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

func printStderr(cmd *exec.Cmd) {
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

func resolveVar(envVar string) string {
	if hasEnvVar(envVar) {
		split := strings.Split(envVar, "/")
		for i, val := range split {
			if hasEnvVar(val) {
				envar := strings.Split(val, "$")
				split[i] = os.Getenv(envar[1])
			}
		}
		return strings.Join(split, "/")
	}
	return envVar
}

func hasEnvVar(str ...string) bool {
	for _, val := range str {
		has := strings.Contains(val, "$")
		if has {
			return has
		}
	}
	return false
}

func beautifyPrint(args []string, dest string) {
	y := color.New(color.FgYellow, color.Bold)
	r := color.New(color.FgHiRed)
	c := color.New(color.FgHiCyan)
	fmt.Println("EXEC SPRUCE:", args[0], args[1], args[2])
	for i := 3; i < len(args); i++ {
		if args[i] == "--prune" {
			r.Printf("\t%s ", args[i])
			i++
			c.Printf("%s \n", args[i])
			continue
		}
		fmt.Printf("\t%s \n", args[i])
	}
	y.Printf("\tto:%s \n\n", dest)
}
