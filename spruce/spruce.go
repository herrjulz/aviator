package spruce

import (
	"regexp"

	yaml "gopkg.in/yaml.v2"

	"github.com/JulzDiverse/aviator"
	"github.com/JulzDiverse/aviator/filemanager"
	"github.com/cppforlife/go-patch/patch"
	"github.com/geofffranks/simpleyaml"
	. "github.com/geofffranks/spruce"
	"github.com/starkandwayne/goutils/ansi"
)

type SpruceClient struct {
	CurlyBraces bool
	store       aviator.FileStore
}

var concourseRegex = `(\{\{|\+\+)([-\_\.\/\w\p{L}\/]+)(\}\}|\+\+)`
var re = regexp.MustCompile("(" + concourseRegex + ")")
var dere = regexp.MustCompile("['\"](" + concourseRegex + ")[\"']")

func New(curlyBraces bool) *SpruceClient {
	return &SpruceClient{
		curlyBraces,
		filemanager.Store(curlyBraces),
	}
}

func NewWithFileFilemanager(filemanager aviator.FileStore, curlyBraces bool) *SpruceClient {
	return &SpruceClient{
		curlyBraces,
		filemanager,
	}
}

func (sc *SpruceClient) MergeWithOpts(options aviator.MergeConf) ([]byte, error) {
	root := make(map[interface{}]interface{})

	err := sc.mergeAllDocs(root, options.Files, options.FallbackAppend, options.EnableGoPatch)
	if err != nil {
		return nil, err
	}

	ev := &Evaluator{Tree: root, SkipEval: options.SkipEval}
	err = ev.Run(options.Prune, options.CherryPicks)
	if err != nil {
		return nil, err
	}

	resultYml, err := yaml.Marshal(ev.Tree)
	if err != nil {
		return nil, err
	}

	return resultYml, nil
}

func (sc *SpruceClient) MergeWithOptsRaw(options aviator.MergeConf) (map[interface{}]interface{}, error) {
	root := make(map[interface{}]interface{})

	err := sc.mergeAllDocs(root, options.Files, options.FallbackAppend, options.EnableGoPatch)
	if err != nil {
		return nil, err
	}

	ev := &Evaluator{Tree: root, SkipEval: options.SkipEval}
	err = ev.Run(options.Prune, options.CherryPicks)

	return ev.Tree, err
}

func (sc *SpruceClient) mergeAllDocs(root map[interface{}]interface{}, paths []string, fallbackAppend bool, goPatchEnabled bool) error {
	m := &Merger{AppendByDefault: fallbackAppend}
	for _, path := range paths {
		var data []byte
		var err error

		data, ok := sc.store.ReadFile(path)
		if !ok {
			return ansi.Errorf("@R{Error reading file from filesystem or internal datastore} @m{%s} \n", path)
		}

		if sc.CurlyBraces {
			data = quoteConcourse(data)
		}

		doc, err := parseYAML(data)
		if err != nil {
			if isArrayError(err) && goPatchEnabled {
				ops, err := parseGoPatch(data)
				if err != nil {
					return ansi.Errorf("@m{%s}: @R{%s}\n", path, err.Error())
				}
				newObj, err := ops.Apply(root)
				if err != nil {
					return ansi.Errorf("@m{%s}: @R{%s}\n", path, err.Error())
				}
				if newRoot, ok := newObj.(map[interface{}]interface{}); !ok {
					return ansi.Errorf("@m{%s}: @R{Unable to convert go-patch output into a hash/map for further merging|\n", path)
				} else {
					root = newRoot
				}
			} else {
				return ansi.Errorf("@m{%s}: @R{%s}\n", path, err.Error())
			}
		} else {
			m.Merge(root, doc)
		}
	}

	return m.Error()
}

func parseYAML(data []byte) (map[interface{}]interface{}, error) {
	y, err := simpleyaml.NewYaml(data)
	if err != nil {
		return nil, err
	}

	doc, err := y.Map()
	if err != nil {
		if _, arrayErr := y.Array(); arrayErr == nil {
			return nil, RootIsArrayError{msg: ansi.Sprintf("@R{Root of YAML document is not a hash/map}: %s\n", err)}
		}
		return nil, ansi.Errorf("@R{Root of YAML document is not a hash/map}: %s\n", err.Error())
	}

	return doc, nil
}

func quoteConcourse(input []byte) []byte {
	return re.ReplaceAll(input, []byte("\"$1\""))
}

func dequoteConcourse(input []byte) string {
	return dere.ReplaceAllString(string(input), "$1")
}

type RootIsArrayError struct {
	msg string
}

func (r RootIsArrayError) Error() string {
	return r.msg
}

func isArrayError(err error) bool {
	_, ok := err.(RootIsArrayError)
	return ok
}

func parseGoPatch(data []byte) (patch.Ops, error) {
	opdefs := []patch.OpDefinition{}
	err := yaml.Unmarshal(data, &opdefs)
	if err != nil {
		return nil, ansi.Errorf("@R{Root of YAML document is not a hash/map. Tried parsing it as go-patch, but got}: %s\n", err)
	}
	ops, err := patch.NewOpsFromDefinitions(opdefs)
	if err != nil {
		return nil, ansi.Errorf("@R{Unable to parse go-patch definitions: %s\n", err)
	}
	return ops, nil
}
