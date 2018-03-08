package filemanager

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/JulzDiverse/mingoak"
	"github.com/starkandwayne/goutils/ansi"
)

type FileManager struct {
	root *mingoak.Dir
}

//var quoteRegexOld = `\{\{([-\_\.\/\w\p{L}\/]+)\}\}`
var quoteRegex = `(\{\{|\+\+)([-\_\.\/\w\p{L}\/]+)(\}\}|\+\+)`
var re = regexp.MustCompile("(" + quoteRegex + ")")
var dere = regexp.MustCompile("['\"](" + quoteRegex + ")[\"']")
var store *FileManager

func Store() *FileManager {
	if store == nil {
		store = &FileManager{mingoak.MkRoot()}
	}
	return store
}

func (ds *FileManager) ReadFile(key string) ([]byte, bool) {
	if _, err := os.Stat(key); os.IsNotExist(err) {
		if re.MatchString(key) {
			key = getKeyFromRegexp(key)
		}
		if file, err := ds.root.ReadFile(key); err == nil {
			return file, true
		}
		return nil, false
	}

	file, err := ioutil.ReadFile(key)
	if err != nil {
		return nil, false
	}
	return file, true
}

func (ds *FileManager) WriteFile(key string, file []byte) error {
	file = dequoteCurlyBraces(file)
	if re.MatchString(key) {
		key = getKeyFromRegexp(key)
		//if _, err := ds.root.ReadFile(key); err == nil {
		//return errors.New(fmt.Sprintf("file %s in virtual filestore already exists", key))
		//}
		ds.root.MkDirAll(getPathFromFilePath(key))
		ds.root.WriteFile(key, []byte(file))
	} else {
		createNonExistingDirs(key)
		err := ioutil.WriteFile(key, file, 0644)
		if err != nil {
			ansi.Errorf("@R{Error writing file} @m{%s}: %s\n", key, err.Error())
		}
	}
	return nil
}

func getPathFromFilePath(filepath string) string {
	sl := strings.Split(filepath, "/")
	sl = sl[:len(sl)-1]
	result := strings.Join(sl, "/")
	return result
}

func (fm *FileManager) ReadDir(path string) ([]os.FileInfo, error) {
	var filePaths []os.FileInfo
	if re.MatchString(path) {

		path = getKeyFromRegexp(path)
		files, err := fm.root.ReadDir(path)
		if err != nil {
			return nil, err
		}
		filePaths = files
	} else {

		files, err := ioutil.ReadDir(path)
		if err != nil {
			return nil, err
		}
		filePaths = files
	}

	return filePaths, nil
}

func (fm *FileManager) Walk(path string) ([]string, error) {
	sl := []string{}
	if re.MatchString(path) {
		path = getKeyFromRegexp(path)
		files, err := fm.root.Walk(path)
		if err != nil {
			return nil, err
		}
		sl = files
	} else {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return nil, err
		} else {
			err := filepath.Walk(path, fillSliceWithFiles(&sl))
			if err != nil {
				return nil, err
			}
		}
	}
	return sl, nil
}

func fillSliceWithFiles(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			*files = append(*files, path)
		}
		return nil
	}
}

func createNonExistingDirs(path string) {
	if strings.Contains(path, "/") {
		sliced := strings.Split(path, "/")
		dirs := sliced[:len(sliced)-1]
		fol := dirs[0]
		for i, dir := range dirs {
			if i > 0 {
				fol = strings.Join([]string{fol, dir}, "/")
			}
			createDir(fol)
		}
	}
}

func createDir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0711)
	}
}

func quoteCurlyBraces(input []byte) []byte {
	return re.ReplaceAll(input, []byte("\"$1\""))
}

func dequoteCurlyBraces(input []byte) []byte {
	return []byte(dere.ReplaceAllString(string(input), "$1"))
}

func getKeyFromRegexp(key string) string {
	matches := re.FindSubmatch([]byte(key))
	return string(matches[len(matches)-2])
}
