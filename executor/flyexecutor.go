package executor

import (
	"fmt"
	"os/exec"
	"reflect"

	"github.com/JulzDiverse/aviator"
	"github.com/pkg/errors"
	"github.com/starkandwayne/goutils/ansi"
)

const (
	setPipelineCmd      = "set-pipeline"
	validatePipelineCmd = "validate-pipeline"
	formatPipelineCmd   = "format-pipeline"
	exposePipelineCmd   = "expose-pipeline"
	hidePipelineCmd     = "hide-pipeline"

	configFlag         = "--config"
	pipelineFlag       = "--pipeline"
	targetFlag         = "--target"
	writeFlag          = "--write"
	strictFlag         = "--strict"
	loadVarsFromFlag   = "--load-vars-from"
	varFlag            = "--var"
	nonInteractiveFlag = "--non-interactive"
)

type FlyExecutor struct{}

func (e FlyExecutor) Command(cfg interface{}) ([]*exec.Cmd, error) {
	fly, ok := cfg.(aviator.Fly)
	if !ok {
		return []*exec.Cmd{}, errors.New(ansi.Sprintf("@R{Type Assertion failed! Cannot assert %s to %s}", reflect.TypeOf(cfg), "aviator.Fly"))
	}

	var args []string
	if fly.ValidatePipeline {
		args = []string{validatePipelineCmd, configFlag, fly.Config}

		if fly.Strict {
			args = append(args, strictFlag)
		}

	} else if fly.FormatPipeline {
		args = []string{formatPipelineCmd, configFlag, fly.Config}

		if fly.Write {
			args = append(args, writeFlag)
		}

	} else {
		args = []string{
			targetFlag, fly.Target, setPipelineCmd, pipelineFlag, fly.Name, configFlag, fly.Config,
		}

		for _, v := range fly.Vars {
			args = append(args, loadVarsFromFlag, v)
		}

		for k, v := range fly.Var {
			args = append(args, varFlag, fmt.Sprintf("%s=%s", k, v))
		}

		if fly.NonInteractive {
			args = append(args, nonInteractiveFlag)
		}
	}

	var exposeArgs []string
	if fly.Expose {
		exposeArgs = []string{targetFlag, fly.Target, exposePipelineCmd, pipelineFlag, fly.Name}
	} else {
		exposeArgs = []string{targetFlag, fly.Target, hidePipelineCmd, pipelineFlag, fly.Name}
	}

	return []*exec.Cmd{
		exec.Command("fly", args...),
		exec.Command("fly", exposeArgs...),
	}, nil
}
