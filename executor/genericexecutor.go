package executor

import (
	"os/exec"
	"reflect"

	"github.com/JulzDiverse/aviator"
	"github.com/pkg/errors"
	"github.com/starkandwayne/goutils/ansi"
)

type GenericExecutor struct{}

func (e GenericExecutor) Command(cfg interface{}) ([]*exec.Cmd, error) {
	execs, ok := cfg.([]aviator.Executable)
	if !ok {
		return []*exec.Cmd{}, errors.New(ansi.Sprintf("@R{Type Assertion failed! Cannot assert %s to %s}", reflect.TypeOf(cfg), "aviator.Exec"))
	}

	cmds := []*exec.Cmd{}
	for _, exe := range execs {
		var args []string
		if len(exe.GlobalOptions) > 0 {
			for _, globOpt := range exe.GlobalOptions {
				args = append(args, globOpt.Name)
				if globOpt.Value != "" {
					args = append(args, globOpt.Value)
				}
			}
		}

		command := exe.Command
		if command.Name != "" {
			args = append(args, command.Name)
			if len(exe.Command.Options) > 0 {
				for _, cmdOpt := range command.Options {
					args = append(args, cmdOpt.Name)
					if cmdOpt.Value != "" {
						args = append(args, cmdOpt.Value)
					}
				}
			}
		}

		if len(exe.Args) > 0 {
			for _, arg := range exe.Args {
				args = append(args, arg)
			}
		}

		cmds = append(cmds, exec.Command(exe.Executable, args...))
	}

	return cmds, nil
}
