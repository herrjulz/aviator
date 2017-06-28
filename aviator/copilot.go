package aviator

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/JulzDiverse/aviator/spruce"
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

func beautifyPrint(opts spruce.MergeOpts, dest string) {
	y := color.New(color.FgYellow, color.Bold)
	r := color.New(color.FgHiRed)
	c := color.New(color.FgHiCyan)
	fmt.Println("SPRUCE MERGE:")
	if len(opts.Prune) != 0 {
		for _, prune := range opts.Prune {
			r.Printf("\t%s ", "--prune")
			c.Printf("  %s \n", prune)
		}
	}
	for _, file := range opts.Files {
		fmt.Printf("\t%s \n", file)
	}
	y.Printf("\tto: %s\n\n", dest)
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func ResolveEnvVars(input []byte) []byte {
	result := os.ExpandEnv(string(input))
	return []byte(result)
}
