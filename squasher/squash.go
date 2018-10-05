package squasher

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
)

func Squash(files [][]byte) ([]byte, error) {
	result := []byte{}
	if len(files) <= 1 {
		return nil, errors.New("zero or one file provided to squash")
	}

	for _, file := range files {
		bytesReader := bytes.NewReader(file)
		bufReader := bufio.NewReader(bytesReader)
		firstLine, _, _ := bufReader.ReadLine()
		if string(firstLine) != "---" && string(firstLine) != "" {
			file = append([]byte(fmt.Sprintf("---\n")), file...)
		}

		file = append(file, []byte(fmt.Sprintf("\n"))...)
		result = append(result, file...)
	}
	return result, nil
}
