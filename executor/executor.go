package executor

import (
	"os"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/starkandwayne/goutils/ansi"
)

func execCmd(cmd *exec.Cmd) error {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, ansi.Sprintf("@R{Failed to run %s}", cmd.Path))
	}

	return nil
}
