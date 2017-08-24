package cockpit

import (
	"regexp"

	"github.com/JulzDiverse/osenv"
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
	Base        string   `yaml:"base"`
	Merge       []Merge  `yaml:"merge"`
	ForEach     ForEach  `yaml:"for_each"`
	Prune       []string `yaml:"prune"`
	CherryPicks []string `yaml:"cherry_pick"`
	SkipEval    bool     `yaml:"skip_eval"`
	To          string   `yaml:"to"`
	ToDir       string   `yaml:"to_dir"`
}

type Merge struct {
	With      With     `yaml:"with"`
	WithIn    string   `yaml:"with_in"`
	WithAllIn string   `yaml:"with_all_in"`
	Except    []string `yaml:"except"`
	Regexp    string   `yaml:"regexp"`
}

type With struct {
	Files []string `yaml:"files"`
	InDir string   `yaml:"in_dir"`
	Skip  bool     `yaml:"skip_non_existing"`
}

type ForEach struct {
	Files          []string `yaml:"files"`
	InDir          string   `yaml:"in_dir"`
	Skip           bool     `yaml:"skip_non_existing"`
	In             string   `yaml:"in"`
	Except         []string `yaml:"except"`
	SubDirs        bool     `yaml:"include_sub_dirs"`
	EnableMatching bool     `yaml:"enable_matching"`
	CopyParents    bool     `yaml:"copy_parents"`
	ForAll         string   `yaml:"for_all"`
	Regexp         string   `yaml:"regexp"`
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
	aviatorYml, err := resolveEnvVars(aviatorYml)
	if err != nil {
		return nil, err
	}

	aviatorYml = quoteCurlyBraces(aviatorYml)
	err = yaml.Unmarshal(aviatorYml, &aviator)
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

func resolveEnvVars(input []byte) ([]byte, error) {
	result, err := osenv.ExpandEnv(string(input))
	return []byte(result), err
}

func quoteCurlyBraces(input []byte) []byte {
	quoteRegex := `\{\{([-\w\p{L}]+)\}\}`
	re := regexp.MustCompile("(" + quoteRegex + ")")
	return re.ReplaceAll(input, []byte("\"$1\""))
}
