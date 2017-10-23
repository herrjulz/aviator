package spruce

import (
	"regexp"

	yaml "gopkg.in/yaml.v2"

	"github.com/JulzDiverse/aviator"
	"github.com/JulzDiverse/aviator/filemanager"
	"github.com/geofffranks/simpleyaml"
	. "github.com/geofffranks/spruce"
	"github.com/starkandwayne/goutils/ansi"
)

type SpruceClient struct {
	store aviator.FileStore
}

var concourseRegex = `(\{\{|\+\+)([-\_\.\/\w\p{L}\/]+)(\}\}|\+\+)`
var re = regexp.MustCompile("(" + concourseRegex + ")")
var dere = regexp.MustCompile("['\"](" + concourseRegex + ")[\"']")

func New() *SpruceClient {
	return &SpruceClient{
		filemanager.Store(),
	}
}

func NewWithFileFilemanager(filemanager aviator.FileStore) *SpruceClient {
	return &SpruceClient{
		filemanager,
	}
}

func (sc *SpruceClient) MergeWithOpts(options aviator.MergeConf) ([]byte, error) {
	root := make(map[interface{}]interface{})

	err := sc.mergeAllDocs(root, options.Files, options.FallbackAppend)
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

	err := sc.mergeAllDocs(root, options.Files, options.FallbackAppend)
	if err != nil {
		return nil, err
	}

	ev := &Evaluator{Tree: root, SkipEval: options.SkipEval}
	err = ev.Run(options.Prune, options.CherryPicks)

	return ev.Tree, err
}

func (sc *SpruceClient) mergeAllDocs(root map[interface{}]interface{}, paths []string, fallbackAppend bool) error {
	m := &Merger{AppendByDefault: fallbackAppend}
	for _, path := range paths {
		var err error

		data, ok := sc.store.ReadFile(path)
		if !ok {
			return ansi.Errorf("@R{Error reading file from filesystem or internal datastore} @m{%s} \n", path)
		}

		data = quoteConcourse(data)

		doc, err := parseYAML(data)
		if err != nil {
			return ansi.Errorf("@m{%s}: @R{%s}\n", path, err.Error())
		}

		err = m.Merge(root, doc)
		if err != nil {
			return err
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
