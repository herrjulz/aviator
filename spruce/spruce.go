package spruce

import (
	"io/ioutil"
	"regexp"

	"github.com/geofffranks/simpleyaml"
	. "github.com/geofffranks/spruce"
	"github.com/starkandwayne/goutils/ansi"
)

type MergeOpts struct {
	Prune          []string
	FallbackAppend bool
	Files          []string
	SkipEval       bool
	CherryPicks    []string
}

var DataStore map[string][]byte = make(map[string][]byte)

var concourseRegex = `\{\{([-\w\p{L}]+)\}\}`

var re = regexp.MustCompile("(" + concourseRegex + ")")

var dere = regexp.MustCompile("['\"](" + concourseRegex + ")[\"']")

func parseYAML(data []byte) (map[interface{}]interface{}, error) {
	y, err := simpleyaml.NewYaml(data)
	if err != nil {
		return nil, err
	}

	doc, err := y.Map()
	if err != nil {
		return nil, ansi.Errorf("@R{Root of YAML document is not a hash/map}: %s\n", err.Error())
	}

	return doc, nil
}

func CmdMergeEval(options MergeOpts) (map[interface{}]interface{}, error) {
	root := make(map[interface{}]interface{})

	err := mergeAllDocs(root, options.Files, options.FallbackAppend)
	if err != nil {
		return nil, err
	}

	ev := &Evaluator{Tree: root, SkipEval: options.SkipEval}
	err = ev.Run(options.Prune, options.CherryPicks)
	return ev.Tree, err
}

func mergeAllDocs(root map[interface{}]interface{}, paths []string, fallbackAppend bool) error {
	m := &Merger{AppendByDefault: fallbackAppend}
	for _, path := range paths {
		var data []byte
		var err error

		data = readYamlFromPathOrStore(path)
		if len(data) == 0 {
			return ansi.Errorf("@R{Error reading file or resolve variable} @m{%s} \n", path)
		}

		data = quoteConcourse(data)

		doc, err := parseYAML(data)
		if err != nil {
			return ansi.Errorf("@m{%s}: @R{%s}\n", path, err.Error())
		}

		m.Merge(root, doc)
	}

	return m.Error()
}

func quoteConcourse(input []byte) []byte {
	return re.ReplaceAll(input, []byte("\"$1\""))
}

func dequoteConcourse(input []byte) string {
	return dere.ReplaceAllString(string(input), "$1")
}

func readYamlFromPathOrStore(path string) []byte {
	var data []byte
	if re.MatchString(path) {
		matches := re.FindSubmatch([]byte(path))
		key := string(matches[len(matches)-1])
		dataTmp, ok := DataStore[key]
		if !ok {
			ansi.Errorf("@R{Error reading variable} @m{%s}\n", key)
			return nil
		}
		data = dataTmp
	} else {
		dataTmp, err := ioutil.ReadFile(path)
		if err != nil {
			ansi.Errorf("@R{Error reading file} @m{%s}: %s\n", path, err.Error())
			return nil
		}
		data = dataTmp
	}
	return data
}

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
