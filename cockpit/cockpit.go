package cockpit

import (

	"github.com/JulzDiverse/aviator"
	"github.com/JulzDiverse/aviator/executor"
	"github.com/JulzDiverse/aviator/processor"
	"github.com/JulzDiverse/aviator/validator"
	"github.com/JulzDiverse/osenv"
	"github.com/pkg/errors"
	"github.com/starkandwayne/goutils/ansi"

	yaml "gopkg.in/yaml.v2"
)

type Cockpit struct {
	spruceProcessor aviator.SpruceProcessor
	flyExecutor     aviator.Executor
	validator       aviator.Validator
}

type Aviator struct {
	cockpit     *Cockpit
	AviatorYaml *aviator.AviatorYaml
}

func Init(
	spruceProcessor aviator.SpruceProcessor,
	flyExecuter aviator.Executor,
	validator aviator.Validator,
) *Cockpit {
	return &Cockpit{spruceProcessor, flyExecuter, validator}
}

func New() *Cockpit {
	return &Cockpit{
		spruceProcessor: processor.New(),
		validator:       validator.New(),
		flyExecutor:     executor.NewFlyExecutor(),
	}
}

func (c *Cockpit) NewAviator(aviatorYml []byte) (*Aviator, error) {
	var aviator aviator.AviatorYaml
	aviatorYml, err := resolveEnvVars(aviatorYml)
	if err != nil {
		return nil, errors.Wrap(err, ansi.Sprintf("@R{Reading Failed}"))
	}

	err = yaml.Unmarshal(aviatorYml, &aviator)
	if err != nil {
		return nil, errors.Wrap(err, ansi.Sprintf("@R{Parsing Failed}"))
	}

	err = c.validator.ValidateSpruce(aviator.Spruce)
	if err != nil {
		return nil, errors.Wrap(err, ansi.Sprintf("@R{Validation Failed}"))
	}

	return &Aviator{c, &aviator}, nil
}

func (a *Aviator) ProcessSprucePlan(verbose bool, silent bool) error {
	err := a.cockpit.spruceProcessor.ProcessWithOpts(a.AviatorYaml.Spruce, verbose, silent)
	if err != nil {
		return errors.Wrap(err, "Processing Spruce Plan FAILED:")
	}
	return nil
}

func (a *Aviator) ExecuteFly() error {
	err := a.cockpit.flyExecutor.Execute(a.AviatorYaml.Fly)
	if err != nil {
		return errors.Wrap(err, "Executing Fly FAILED")
	}
	return nil
}

func resolveEnvVars(input []byte) ([]byte, error) {
	result, err := osenv.ExpandEnv(string(input))
	return []byte(result), err
}
