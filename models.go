package aviator

import "os"

type AviatorYaml struct {
	Spruce []Spruce `yaml:"spruce"`
	Squash Squash   `yaml:"squash"`
	Fly    Fly      `yaml:"fly"`
	Kube   Kube     `yaml:"kubectl"`
}

type Spruce struct {
	Base        string   `yaml:"base"`
	Merge       []Merge  `yaml:"merge"`
	ForEach     ForEach  `yaml:"for_each"`
	Prune       []string `yaml:"prune"`
	CherryPicks []string `yaml:"cherry_pick"`
	SkipEval    bool     `yaml:"skip_eval"`
	GoPatch     bool     `yaml:"go_patch"`
	To          string   `yaml:"to"`
	ToDir       string   `yaml:"to_dir"`
	Modify      Modify   `yaml:"modify"`
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
	Name           string            `yaml:"name"`
	Target         string            `yaml:"target"`
	Config         string            `yaml:"config"`
	Vars           []string          `yaml:"load_vars_from"`
	Expose         bool              `yaml:"expose"`
	Var            map[string]string `yaml:"vars"`
	NonInteractive bool              `yaml:"non_interactive"`
}

type Kube struct {
	Apply KubeApply `yaml:"apply"`
}

type KubeApply struct {
	File      string `yaml:"file"`
	Force     bool   `yaml:"force"`
	DryRun    bool   `yaml:"dry_run"`
	Overwrite bool   `yaml:"no_overwrite"`
	Recursive bool   `yaml:"recursive"`
	Output    string `yaml:"output"`
}

type MergeConf struct {
	Files          []string
	Prune          []string
	CherryPicks    []string
	SkipEval       bool
	FallbackAppend bool
	EnableGoPatch  bool
}

type Modify struct {
	Delete []string  `yaml:"delete"`
	Set    []PathVal `yaml:"set"`
	Update []PathVal `yaml:"update"`
}

type PathVal struct {
	Path  string `yaml:"path"`
	Value string `yaml:"value"`
}

type Squash struct {
	Contents []SquashContent `yaml:"contents"`
	To       string          `yaml:"to"`
}

type SquashContent struct {
	Files  []string `yaml:"files"`
	Except []string `yaml:"except"`
	Dir    string   `yaml:"dir"`
}

//go:generate counterfeiter . SpruceProcessor
type SpruceProcessor interface {
	Process([]Spruce) error
	ProcessWithOpts([]Spruce, bool, bool) error
}

//go:generate counterfeiter . FlyExecuter
type FlyExecuter interface {
	Execute(interface{}) error
}

//go:generate counterfeiter . Executor
type Executor interface {
	Execute(interface{}) error
}

//go:generate counterfeiter . SpruceClient
type SpruceClient interface {
	MergeWithOpts(MergeConf) ([]byte, error)
}

//go:generate counterfeiter . FileStore
type FileStore interface {
	ReadFile(string) ([]byte, bool)
	WriteFile(string, []byte) error
	ReadDir(string) ([]os.FileInfo, error)
	Walk(string) ([]string, error)
}

//go:generate counterfeiter . Validator
type Validator interface {
	ValidateSpruce([]Spruce) error
}

//go:generate counterfeiter . Modifier
type Modifier interface {
	Modify([]byte, Modify) ([]byte, error)
}

//go:generate counterfeiter . GomlClient
type GomlClient interface {
	Delete([]byte, string) ([]byte, error)
	Set([]byte, string, string) ([]byte, error)
	Update([]byte, string, string) ([]byte, error)
}
