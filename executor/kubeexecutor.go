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
	kustomizeFlag = "--kustomize"
	forceFlag     = "--force"
	filenameFlag  = "--filename"
	dryRunFlag    = "--dry-run"
	overwriteFlag = "--overwrite"
	validateFlag  = "--validate"
	outputFlag    = "--output"
	recursiveFlag = "--recursive"
)

type KubeExecutor struct{}

func (e KubeExecutor) Command(cfg interface{}) ([]*exec.Cmd, error) {
	kube, ok := cfg.(aviator.Kube)
	if !ok {
		return []*exec.Cmd{}, errors.New(ansi.Sprintf("@R{Type Assertion failed! Cannot assert %s to %s}", reflect.TypeOf(cfg), "aviator.Kube"))
	}

	apply := kube.Apply

	var args []string
	if apply.Kustomize {
		args = []string{
			"apply", kustomizeFlag, apply.File,
		}
	} else {
		args = []string{
			"apply", filenameFlag, apply.File,
		}
	}

	if apply.Recursive {
		args = append(args, fmt.Sprintf("%s=%s", recursiveFlag, "true"))
	}

	if apply.Force {
		args = append(args, forceFlag)
	}

	if apply.DryRun {
		args = append(args, dryRunFlag)
	}

	if apply.Overwrite {
		args = append(args, overwriteFlag)
	}

	if apply.Validate {
		args = append(args, validateFlag)
	}

	if apply.Output != "" {
		args = append(args, outputFlag, apply.Output)
	}

	return []*exec.Cmd{exec.Command("kubectl", args...)}, nil
}
