package main

import (
	"fmt"
	"io/ioutil"
	"masterjulz/aviator/aviator"
	"os"

	"github.com/urfave/cli"
)

func main() {

	cmd := setCli()

	cmd.Action = func(c *cli.Context) error {
		if c.String("t") == "" {
			fmt.Println("Please specify target")
			os.Exit(1)
		}

		if c.String("p") == "" {
			fmt.Println("Please specify a pipline name")
			os.Exit(1)
		}

		target := c.String("target")

		pipeline := c.String("pipeline")

		fmt.Println("Target set to", target)

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
			aviator.ProcessSpruceChain(yml.Spruce)
			aviator.FlyPipeline(yml.Fly, target, pipeline)
		}

		return nil
	}
	cmd.Run(os.Args)
}
