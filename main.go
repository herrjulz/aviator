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
		target := ""
		pipeline := ""

		if c.String("t") != "" {
			target = c.String("target")
		}

		if c.String("p") != "" {
			pipeline = c.String("pipeline")
		}

		aviatorFile := "./aviator.yml"

		var yml aviator.Aviator
		if _, err := os.Stat(aviatorFile); os.IsNotExist(err) {
			fmt.Println("No Aviator file found. Please navigate to a Aviator directory and run Aviator again")
		} else {
			ymlBytes, err := ioutil.ReadFile(aviatorFile)
			if err != nil {
				panic(err)
			}
			yml = aviator.ReadYaml(ymlBytes)

			err := aviator.ProcessSprucePlan(yml.Spruce)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			if target != "" {
				fmt.Println("Target set to", target)
				aviator.FlyPipeline(yml.Fly, target, pipeline)
			}

		}

		return nil
	}
	cmd.Run(os.Args)
}
