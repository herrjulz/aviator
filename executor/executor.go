package executor

import (
	"os"
	"os/exec"

	"code.cloudfoundry.org/commandrunner"

	"github.com/pkg/errors"
	"github.com/starkandwayne/goutils/ansi"
)

func execCmd(cmdname string, args []string, runner commandrunner.CommandRunner) error {
	cmd := exec.Command(cmdname, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := runner.Run(cmd)
	if err != nil {
		return errors.Wrap(err, ansi.Sprintf("@R{Failed to run %s}", cmdname))
	}

	return nil
}
