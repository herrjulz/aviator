package main

import "github.com/urfave/cli"

func setCli() *cli.App {
	cmd := cli.NewApp()
	cmd.Authors = []cli.Author{
		cli.Author{
			Name:  "Julz Skupnjak",
			Email: "julian.skupnjak@gmail.com",
		},
	}
	cmd.Name = "Aviator"
	cmd.Usage = "CLI Tool to Run AVIATOR Concourse Pipelines"
	cmd.Version = "0.0.1"
	cmd.Flags = getFlags()
	return cmd
}

func getFlags() []cli.Flag {
	var flags []cli.Flag
	flags = []cli.Flag{}
	return flags
}
