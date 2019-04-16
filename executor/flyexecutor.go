package executor

import (
	"fmt"
	"os/exec"
	"reflect"

	"github.com/JulzDiverse/aviator"
	"github.com/pkg/errors"
	"github.com/starkandwayne/goutils/ansi"
)

type FlyExecutor struct{}

func (e FlyExecutor) Command(cfg interface{}) (*exec.Cmd, error) {
	fly, ok := cfg.(aviator.Fly)
	if !ok {
		return &exec.Cmd{}, errors.New(ansi.Sprintf("@R{Type Assertion failed! Cannot assert %s to %s}", reflect.TypeOf(cfg), "aviator.Fly"))
	}

	var args []string
	if fly.ValidatePipeline {
		args = []string{"validate-pipeline", "-c", fly.Config}

		if fly.Strict {
			args = append(args, "--strict")
		}

	} else if fly.FormatPipeline {
		args = []string{"format-pipeline", "-c", fly.Config}

		if fly.Write {
			args = append(args, "--write")
		}

	} else {
		args = []string{
			"-t", fly.Target, "set-pipeline", "-p", fly.Name, "-c", fly.Config,
		}

		for _, v := range fly.Vars {
			args = append(args, "-l", v)
		}

		for k, v := range fly.Var {
			args = append(args, "-v", fmt.Sprintf("%s=%s", k, v))
		}

		if fly.NonInteractive {
			args = append(args, "-n")
		}
	}

	return exec.Command("fly", args...), nil
}

func (e FlyExecutor) Execute(cmd *exec.Cmd, cfg interface{}) error {
	fly, ok := cfg.(aviator.Fly)
	if !ok {
		return errors.New(ansi.Sprintf("@R{Type Assertion failed! Cannot assert %s to %s}", reflect.TypeOf(cfg), "aviator.Fly"))
	}

	err := execCmd(cmd)
	if err != nil {
		return err
	}

	if fly.Expose {
		args := []string{"-t", fly.Target, "expose-pipeline", "-p", fly.Name}
		cmd = exec.Command("fly", args...)
		err := execCmd(cmd)
		if err != nil {
			return err
		}
	} else {
		args := []string{"-t", fly.Target, "hide-pipeline", "-p", fly.Name}
		cmd = exec.Command("fly", args...)
		err := execCmd(cmd)
		if err != nil {
			return err
		}
	}

	return nil
}
