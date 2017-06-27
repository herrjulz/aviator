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
}

var concourseRegex = `\{\{([-\w\p{L}]+)\}\}`

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

	ev := &Evaluator{Tree: root, SkipEval: false}
	err = ev.Run(options.Prune, []string{})
	return ev.Tree, err
}

func mergeAllDocs(root map[interface{}]interface{}, paths []string, fallbackAppend bool) error {
	m := &Merger{AppendByDefault: fallbackAppend}
	for _, path := range paths {
		var data []byte
		var err error

		data, err = ioutil.ReadFile(path)
		if err != nil {
			return ansi.Errorf("@R{Error reading file} @m{%s}: %s\n", path, err.Error())
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
	re := regexp.MustCompile("(" + concourseRegex + ")")
	return re.ReplaceAll(input, []byte("\"$1\""))
}
