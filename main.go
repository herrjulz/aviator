package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/JulzDiverse/aviator/aviator"

	"github.com/urfave/cli"
)

func main() {

	cmd := setCli()

	cmd.Action = func(c *cli.Context) error {
		aviatorFile := "./aviator.yml"

		var yml aviator.Aviator
		if _, err := os.Stat(aviatorFile); os.IsNotExist(err) {
			fmt.Println("No Aviator file found. Please navigate to a Aviator directory and run Aviator again")
		} else {
			ymlBytes, err := ioutil.ReadFile(aviatorFile)
			if err != nil {
				panic(err)
			}

			yml = aviator.ReadYaml(aviator.ResolveEnvVars(ymlBytes))
			err = aviator.ProcessSprucePlan(yml.Spruce)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			if yml.Fly.Target != "" && yml.Fly.Name != "" && yml.Fly.Config != "" {
				fmt.Println("Target set to", yml.Fly.Target)
				aviator.FlyPipeline(yml.Fly)
			}

		}

		return nil
	}
	cmd.Run(os.Args)
}
