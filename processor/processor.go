package processor

import (
	"fmt"

	yaml "gopkg.in/yaml.v2"
)

type Processor struct {
	Aviator Aviator
}

type Aviator struct {
	Spruce []Spruce `yaml:"spruce"`
}

type Spruce struct {
	Base  string  `yaml:"base"`
	Merge []Merge `yaml:"merge"`
	To    string  `yaml:"to"`
}

type Merge struct {
	With With `yaml:"with"`
}

type With struct {
	Files []string `yaml:"files"`
}

func New(aviatorYml []byte) (*Processor, error) {
	var aviator Aviator
	err := yaml.Unmarshal(aviatorYml, &aviator)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%s", aviator)
	return &Processor{aviator}, nil
}
