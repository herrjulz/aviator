package cockpit

import (
	"os"
	"regexp"

	"github.com/pkg/errors"

	yaml "gopkg.in/yaml.v2"
)

type Cockpit struct {
	spruceProcessor SpruceProcessor
	flyExecuter     FlyExecuter
}

type Aviator struct {
	cockpit     *Cockpit
	AviatorYaml *AviatorYaml
}

type AviatorYaml struct {
	Spruce []Spruce `yaml:"spruce"`
	Fly    Fly      `yaml:"fly"`
}

type Spruce struct {
	Base           string   `yaml:"base"`
	Merge          []Merge  `yaml:"merge"`
	To             string   `yaml:"to"`
	ToDir          string   `yaml:"to_dir"`
	ForEach        []string `yaml:"for_each"`
	ForEachIn      string   `yaml:"for_each_in"`
	Except         []string `yaml:"except"`
	WalkThrough    string   `yaml:"walk_through"`
	Prune          []string `yaml:"prune"`
	CherryPicks    []string `yaml:"cherry_pick"`
	EnableMatching bool     `yaml:"enable_matching"`
	CopyParents    bool     `yaml:"copy_parents"`
	SkipEval       bool     `yaml:"skip_eval"`
	ForAll         string   `yaml:"for_all"`
}

type Merge struct {
	With   With     `yaml:"with"`
	WithIn string   `yaml:"with_in"`
	Except []string `yaml:"except"`
	Regexp string   `yaml:"regexp"`
	Skip   bool     `yaml:"skip_non_existing"`
}

type With struct {
	Files []string `yaml:"files"`
}

type Fly struct {
	Name   string   `yaml:"name"`
	Target string   `yaml:"target"`
	Config string   `yaml:"config"`
	Vars   []string `yaml:"vars"`
	Expose bool     `yaml:"expose"`
}

//go:generate counterfeiter . SpruceProcessor
type SpruceProcessor interface {
	Process([]Spruce) ([]byte, error)
}

//go:generate counterfeiter . FlyExecuter
type FlyExecuter interface {
	Execute(Fly) error
}

func Init(
	spruceProcessor SpruceProcessor,
	flyExecuter FlyExecuter,
) *Cockpit {
	return &Cockpit{spruceProcessor, flyExecuter}
}

func (c *Cockpit) NewAviator(aviatorYml []byte) (*Aviator, error) {
	var aviator AviatorYaml
	aviatorYml = resolveEnvVars(aviatorYml)
	aviatorYml = quoteCurlyBraces(aviatorYml)
	err := yaml.Unmarshal(aviatorYml, &aviator)
	if err != nil {
		return nil, err
	}
	return &Aviator{c, &aviator}, nil
}

func (a *Aviator) ProcessSprucePlan() ([]byte, error) {
	mergedYaml, err := a.cockpit.spruceProcessor.Process(a.AviatorYaml.Spruce)
	if err != nil {
		return nil, errors.Wrap(err, "Processing Spruce Plan FAILED:")
	}
	return mergedYaml, nil
}

func (a *Aviator) ExecuteFly() error {
	err := a.cockpit.flyExecuter.Execute(a.AviatorYaml.Fly)
	if err != nil {
		return errors.Wrap(err, "Executing Fly FAILED")
	}
	return nil
}

func resolveEnvVars(input []byte) []byte {
	result := os.ExpandEnv(string(input))
	return []byte(result)
}

func quoteCurlyBraces(input []byte) []byte {
	quoteRegex := `\{\{([-\w\p{L}]+)\}\}`
	re := regexp.MustCompile("(" + quoteRegex + ")")
	return re.ReplaceAll(input, []byte("\"$1\""))
}
