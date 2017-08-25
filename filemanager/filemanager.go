package filemanager

import (
	"errors"
	"fmt"
	"regexp"
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
	if re.MatchString(key) {
		matches := re.FindSubmatch([]byte(key))
		key = string(matches[len(matches)-1])
	}

	if file, ok := ds.files[key]; ok {
		return file, true
	}

	return nil, false
}

func (ds *FileStore) WriteFile(key string, file []byte) error {
	if re.MatchString(key) {
		matches := re.FindSubmatch([]byte(key))
		key = string(matches[len(matches)-1])
	}

	if _, ok := ds.files[key]; ok {
		return errors.New(fmt.Sprintf("file %s in virtual filestore already exists", key))
	}

	file = dequoteCurlyBraces(file)
	ds.files[key] = []byte(file)
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
