package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/JulzDiverse/aviator/cmd/aviator/cockpit"
	"github.com/JulzDiverse/aviator/validator"
	"github.com/starkandwayne/goutils/ansi"
	"github.com/urfave/cli"
)

func main() {
	cmd := setCli()

	cmd.Action = func(c *cli.Context) error {
		aviatorFile := c.String("file")
		if !verifyAviatorFileExists(aviatorFile) {
			exitWithNoAviatorFile()
		} else {
			vars := c.StringSlice("var")
			varsMap := varsToMap(vars)

			aviatorYml, err := ioutil.ReadFile(aviatorFile)
			exitWithError(err)

			cockpit := cockpit.New(
				c.Bool("curly-braces"),
				c.Bool("dry-run"),
			)

			aviator, err := cockpit.NewAviator(
				aviatorYml,
				varsMap,
				c.Bool("silent"),
				c.Bool("verbose"),
				c.Bool("dry-run"),
			)

			handleError(err)

			err = aviator.ProcessSprucePlan()
			exitWithError(err)

			squash := aviator.AviatorYaml.Squash
			if len(squash.Contents) != 0 {
				err = aviator.ProcessSquashPlan()
				exitWithError(err)
			}

			if !c.Bool("dry-run") {
				fly := aviator.AviatorYaml.Fly
				if fly.Name != "" && fly.Target != "" && fly.Config != "" {
					err = aviator.ExecuteFly()
					exitWithError(err)
				}

				kube := aviator.AviatorYaml.Kube.Apply
				if kube.File != "" {
					err = aviator.ExecuteKube()
					exitWithError(err)
				}

				exec := aviator.AviatorYaml.Exec
				if len(exec) != 0 {
					err = aviator.ExecuteGeneric()
					exitWithError(err)
				}
			}
		}

		return nil
	}
	cmd.Run(os.Args)
}

func varsToMap(vars []string) map[string]string {
	result := map[string]string{}
	for _, v := range vars {
		sl := strings.Split(v, "=")
		result[sl[0]] = sl[1]
	}
	return result
}

func verifyAviatorFileExists(file string) bool {
	if file == "aviator.yml" {
		if _, err := os.Stat(file); !os.IsNotExist(err) {
			return true
		}
	} else {
		if _, err := os.Stat(file); !os.IsNotExist(err) {
			return true
		}
	}
	return false
}

func exitWithNoAviatorFile() {
	ansi.Printf("@R{No Aviator file found.}\n\n")
	fmt.Println("Please navigate to a directory that contains an aviator.yml or specify a AVIATOR YAML with [--file|-f] option and run aviator again")
	os.Exit(1)
}

func exitWithError(err error) {
	if err != nil {
		ansi.Printf("@R{%s}\n", err.Error())
		os.Exit(1)
	}
}

func handleError(err error) {
	if err != nil {
		switch err.(type) {
		case validator.MergeCombinationError:
			printMergeCombinationError(err)
		case validator.MergeWithCombinationError:
			printMergeWithCombinationError(err)
		case validator.MergeRegexpCombinationError:
			printMergeRegexpCombinationError(err)
		case validator.MergeExceptCombinationError:
			printMergeExceptCombinationError(err)
		case validator.ForEachCombinationError:
			printForEachCombinationError(err)
		case validator.ForEachFilesCombinationError:
			printForEachFilesCombinationError(err)
		case validator.ForEachInCombinationError:
			printForEachInCombinationError(err)
		case validator.ForEachWalkCombinationError:
			printForEachWalkCombinationError(err)
		case validator.ForEachRegexpCombinationError:
			printForEachRegexpCombinationError(err)
		default:
			ansi.Printf(err.Error())
		}
		os.Exit(1)
	}
}
