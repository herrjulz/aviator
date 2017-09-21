package mingoak

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func MkRoot() *Dir {
	return &Dir{
		components: map[string]os.FileInfo{},
		name:       "root",
		time:       time.Now(),
	}
}

func (d *Dir) WriteFile(path string, file []byte) error {
	if path == "" {
		return errors.New("No file name or path provided!")
	}

	current := d
	sl := slicePath(path)
	for i, name := range sl {
		if i == len(sl)-1 {
			fileInfo := File{
				content: file,
				name:    name,
				time:    time.Now(),
			}
			current.components[name] = &fileInfo
			current.componentsl = append(current.componentsl, &fileInfo)
			break
		}
		current = current.components[name].(*Dir)
	}
	return nil
}

func (d *Dir) ReadFile(path string) ([]byte, error) {
	current := d
	for _, name := range slicePath(path) {
		if result, ok := current.components[name]; ok && result.IsDir() {
			current = result.(*Dir)
		} else if result, ok := current.components[name]; ok {
			file := result.(*File)
			return file.content, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("File %s not found!", path))
}

func (d *Dir) MkDirAll(path string) {
	current := d
	for _, name := range slicePath(path) {
		_, ok := current.components[name]
		if !ok {
			dir := Dir{
				components: map[string]os.FileInfo{},
				name:       name,
				time:       time.Now(),
			}
			current.components[name] = &dir
			current.componentsl = append(current.componentsl, &dir)
		}
		current = current.components[name].(*Dir)
	}
}

func (d *Dir) ReadDir(dirname string) ([]os.FileInfo, error) {
	dir, err := d.getDir(dirname)
	if err != nil {
		return nil, err
	}
	return dir.componentsl, nil
}

func (d *Dir) Walk(path string) ([]string, error) {
	dir, err := d.getDir(path)
	if err != nil {
		return nil, err
	}

	files := walkRecursion(dir, path)
	return files, nil
}

func walkRecursion(dir *Dir, basepath string) []string {
	files := []string{}
	for k, v := range dir.components {
		if v.IsDir() {
			subFiles := walkRecursion(v.(*Dir), filepath.Join(basepath, k))
			files = append(files, subFiles...)
		} else {
			files = append(files, filepath.Join(basepath, k))
		}
	}
	sort.Strings(files)
	return files
}

func (d *Dir) getDir(path string) (*Dir, error) {
	current := d
	for _, name := range slicePath(path) {
		if result, ok := current.components[name]; ok && result.IsDir() {
			current = result.(*Dir)
		} else {
			return &Dir{}, errors.New(fmt.Sprintf("Directory %s not found!", path))
		}
	}
	return current, nil
}

func slicePath(path string) []string {
	sl := strings.Split(path, "/")
	if sl[len(sl)-1] == "" {
		sl = sl[:len(sl)-1]
	}
	return sl
}
