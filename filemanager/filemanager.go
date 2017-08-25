package filemanager

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/starkandwayne/goutils/ansi"
)

type FileStore struct {
	files map[string][]byte
}

var quoteRegex = `\{\{([-\w\p{L}]+)\}\}`
var re = regexp.MustCompile("(" + quoteRegex + ")")
var dere = regexp.MustCompile("['\"](" + quoteRegex + ")[\"']")
var store *FileStore

func Store() *FileStore {
	if store == nil {
		store = &FileStore{map[string][]byte{}}
	}
	return store
}

func (ds *FileStore) ReadFile(key string) ([]byte, bool) {
	if _, err := os.Stat(key); os.IsNotExist(err) {
		if re.MatchString(key) {
			matches := re.FindSubmatch([]byte(key))
			key = string(matches[len(matches)-1])
		}
		if file, ok := ds.files[key]; ok {
			return file, true
		}
	}

	file, err := ioutil.ReadFile(key)
	if err != nil {
		return nil, false
	}

	return file, true
}

func (ds *FileStore) WriteFile(key string, file []byte) error {
	file = dequoteCurlyBraces(file)
	if re.MatchString(key) {
		matches := re.FindSubmatch([]byte(key))
		key = string(matches[len(matches)-1])

		if _, ok := ds.files[key]; ok {
			return errors.New(fmt.Sprintf("file %s in virtual filestore already exists", key))
		}
		ds.files[key] = []byte(file)
	} else {
		err := ioutil.WriteFile(key, file, 0644)
		if err != nil {
			ansi.Errorf("@R{Error writing file} @m{%s}: %s\n", key, err.Error())
		}
	}

	return nil
}

func quoteCurlyBraces(input []byte) []byte {
	return re.ReplaceAll(input, []byte("\"$1\""))
}

func dequoteCurlyBraces(input []byte) []byte {
	return []byte(dere.ReplaceAllString(string(input), "$1"))
}

func (ds *FileStore) PrintFiles() {
	for key, file := range ds.files {
		fmt.Println(key, string(file))
	}
}
