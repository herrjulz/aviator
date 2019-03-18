package executor

import (
	"reflect"

	"code.cloudfoundry.org/commandrunner"
	"code.cloudfoundry.org/commandrunner/linux_command_runner"
	"github.com/JulzDiverse/aviator"
	"github.com/pkg/errors"
	"github.com/starkandwayne/goutils/ansi"
)

type KubeExecutor struct {
	runner commandrunner.CommandRunner
}

func NewKubeExecutorWithCustomRunner(runner commandrunner.CommandRunner) *KubeExecutor {
	return &KubeExecutor{
		runner,
	}
}

func NewKubeExecutor() *KubeExecutor {
	return &KubeExecutor{
		runner: linux_command_runner.New(),
		//runner: windows_command_runner.New(false),
	}
}

func (e *KubeExecutor) Execute(cfg interface{}) error {
	kube, ok := cfg.(aviator.Kube)
	if !ok {
		return errors.New(ansi.Sprintf("@R{Type Assertion failed! Cannot assert %s to %s}", reflect.TypeOf(cfg), "aviator.Kube"))
	}

	apply := kube.Apply

	args := []string{
		"apply", "-f", apply.File,
	}

	if apply.Force {
		args = append(args, "--force")
	}

	if apply.DryRun {
		args = append(args, "--dry-run")
	}

	if apply.Overwrite {
		args = append(args, "--overwrite")
	}

	if apply.Recursive {
		args = append(args, "--recursive")
	}

	if apply.Output != "" {
		args = append(args, "--output", apply.Output)
	}

	return execCmd("kubectl", args, e.runner)
}
