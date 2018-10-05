package processor

import (
	"path/filepath"
	"regexp"

	"github.com/JulzDiverse/aviator"
)

type FileProcessor struct {
	Store aviator.FileStore
}

func (f *FileProcessor) CollectFilesFromDir(dir, regex string, ignore []string) []string {
	result := []string{}
	if dir != "" {
		files, _ := f.Store.ReadDir(dir)

		for _, f := range files {
			if except(ignore, f.Name()) {
				continue
			}

			matched, _ := regexp.MatchString(regex, f.Name())
			if !f.IsDir() && matched {
				result = append(result, filepath.Join(resolveBraces(dir)+f.Name()))
			}
		}
	}
	return result
}
