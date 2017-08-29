package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/JulzDiverse/aviator/cockpit"
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

			aviatorYml, err := ioutil.ReadFile(aviatorFile)
			exitWithError(err)

			cockpit := cockpit.New()
			aviator, err := cockpit.NewAviator(aviatorYml)
			handleError(err)

			err = aviator.ProcessSprucePlan(c.Bool("verbose"), c.Bool("silent"))
			exitWithError(err)

			fly := aviator.AviatorYaml.Fly
			if fly.Name != "" && fly.Target != "" && fly.Config != "" {
				aviator.ExecuteFly()
			}
		}

		return nil
	}
	cmd.Run(os.Args)
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
		fmt.Println(err.Error())
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
