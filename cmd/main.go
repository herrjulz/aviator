package main

import "os"

func WriteYamlToPathOrStore(path string, data []byte) {
	if re.MatchString(path) {
		dataString := dequoteConcourse(data)
		matches := re.FindSubmatch([]byte(path))
		key := string(matches[len(matches)-1])
		DataStore[key] = []byte(dataString)
	} else {
		dataString := dequoteConcourse(data)
		err := ioutil.WriteFile(path, []byte(dataString), 0644)
		if err != nil {
			ansi.Errorf("@R{Error writing file} @m{%s}: %s\n", path, err.Error())
		}
	}
}

func CreateDir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0711)
	}
}
