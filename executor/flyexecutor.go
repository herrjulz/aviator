package executor

import (
	"fmt"
	"reflect"

	"code.cloudfoundry.org/commandrunner"
	"code.cloudfoundry.org/commandrunner/linux_command_runner"
	"github.com/JulzDiverse/aviator"
	"github.com/pkg/errors"
	"github.com/starkandwayne/goutils/ansi"
)

type FlyExecutor struct {
	runner commandrunner.CommandRunner
}

func NewFlyExecutorWithCustomRunner(runner commandrunner.CommandRunner) *FlyExecutor {
	return &FlyExecutor{
		runner,
	}
}

func NewFlyExecutor() *FlyExecutor {
	return &FlyExecutor{
		runner: linux_command_runner.New(),
	}
}

func (e *FlyExecutor) Execute(cfg interface{}) error {
	fly, ok := cfg.(aviator.Fly)
	if !ok {
		return errors.New(ansi.Sprintf("@R{Type Assertion failed! Cannot assert %s to %s}", reflect.TypeOf(cfg), "aviator.Fly"))
	}

	args := []string{
		"-t", fly.Target, "set-pipeline", "-p", fly.Name, "-c", fly.Config,
	}

	for _, v := range fly.Vars {
		args = append(args, "-l", v)
	}

	for k, v := range fly.Var {
		args = append(args, "-v", fmt.Sprintf("%s=%s", k, v))
	}

	if fly.NonInteractive == true {
		args = append(args, "-n")
	}

	err := execCmd("fly", args)
	if err != nil {
		return errors.Wrap(err, ansi.Sprintf("@R{Failed to execute Fly}"))
	}

	if fly.Expose {
		args = []string{"-t", fly.Target, "expose-pipeline", "-p", fly.Name}
		err := execCmd("fly", args)
		if err != nil {
			return errors.Wrap(err, ansi.Sprintf("@R{Failed to execute Fly}"))
		}
	}
	return nil
}

func (e *FlyExecutor) ExecuteWithCustomRunner(cfg interface{}) error {
	fly, ok := cfg.(aviator.Fly)
	if !ok {
		return errors.New(ansi.Sprintf("@R{Type Assertion failed! Cannot assert %s to %s}", reflect.TypeOf(cfg), "aviator.Fly"))
	}

	args := []string{
		"-t", fly.Target, "set-pipeline", "-p", fly.Name, "-c", fly.Config,
	}

	for _, v := range fly.Vars {
		args = append(args, "-l", v)
	}

	err := execCmdWithRunner("fly", args, e.runner)
	if err != nil {
		return errors.Wrap(err, ansi.Sprintf("@R{Failed to execute Fly}"))
	}

	if fly.Expose {
		args = []string{"-t", fly.Target, "expose-pipeline", "-p", fly.Name}
		err := execCmdWithRunner("fly", args, e.runner)
		if err != nil {
			return errors.Wrap(err, ansi.Sprintf("@R{Failed to execute Fly}"))
		}
	}
	return nil
}
