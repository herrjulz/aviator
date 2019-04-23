package executor

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/starkandwayne/goutils/ansi"
)

type Executor struct {
	silent bool
}

func New(silent bool) *Executor {
	return &Executor{
		silent: silent,
	}
}

func (e *Executor) Execute(cmds []*exec.Cmd) error {
	for _, c := range cmds {
		if !e.silent {
			fmt.Println(stringifyCmd(c))
		}
		err := e.execCmd(c)
		if err != nil {
			return err
		}
		if !e.silent {
			fmt.Println("")
		}
	}

	return nil
}

func (e *Executor) execCmd(cmd *exec.Cmd) error {
	if !e.silent {
		cmd.Stdout = os.Stdout
	}
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, ansi.Sprintf("@R{Failed to run %s}", cmd.Path))
	}

	return nil
}

func stringifyCmd(cmd *exec.Cmd) string {
	result := ""
	result = ansi.Sprintf("@G{AVIATOR EXECUTE:$} %s", cmd.Args[0])
	for i := 1; i < len(cmd.Args); i++ {
		result = fmt.Sprintf("%s %s", result, cmd.Args[i])
	}
	return result
}
