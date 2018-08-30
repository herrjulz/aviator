package evaluator

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var variableFormatRegex = regexp.MustCompile(`\(\(\s([-\w\p{L}]+)\s\)\)`)

func Evaluate(aviatorFile []byte, vars map[string]string) ([]byte, error) {
	var err error
	return variableFormatRegex.ReplaceAllFunc(aviatorFile, func(match []byte) []byte {
		key := string(variableFormatRegex.FindSubmatch(match)[1])

		val, ok := vars[key]
		if !ok {
			err = errors.New(fmt.Sprintf("Variable (( %s )) not provided", key))
		}

		var replace []byte
		if strings.Contains(val, "\n") {
			replace, _ = json.Marshal(val)
		} else {
			replace = []byte(val)
		}

		return []byte(replace)
	}), err
}
