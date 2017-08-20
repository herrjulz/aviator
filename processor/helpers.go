package processor

import (
	"log"
	"os"
	"path/filepath"

	"github.com/JulzDiverse/aviator/cockpit"
)

func except(except []string, file string) bool {
	for _, f := range except {
		if f == file {
			return true
		}
	}
	return false
}

func getRegexp(merge cockpit.Merge) string {
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

func mergeType(cfg cockpit.Spruce) string {
	if (cfg.ForEach.Files == nil ||
		len(cfg.ForEach.Files) > 0) &&
		cfg.ForEach.In == "" {
		return "default"
	}
	if len(cfg.ForEach.Files) > 0 {
		return "forEach"
	}
	if cfg.ForEach.In != "" && cfg.ForEach.SubDirs == false {
		return "forEachIn"
	}
	if cfg.ForEach.In != "" && cfg.ForEach.SubDirs == true {
		if cfg.ForEach.ForAll != "" {
			return "walkThrough"
		} else {
			return "walkThroughForAll"
		}
	}
	return ""
}

func getAllFilesIncludingSubDirs(path string) []string {
	sl := []string{}
	err := filepath.Walk(path, fillSliceWithFiles(&sl))
	if err != nil {
		log.Fatal(err)
	}
	return sl
}

func fillSliceWithFiles(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			*files = append(*files, path)
		}
		return nil
	}
}

//func concatFileNameWithPath(path string) (string, string) {
//chunked := strings.Split(path, "/")
//fileName := chunked[len(chunked)-2] + "_" + chunked[len(chunked)-1]
//parent := chunked[len(chunked)-2]
//return fileName, parent
//}
