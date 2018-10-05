package main

import "github.com/urfave/cli"

func setCli() *cli.App {
	cmd := cli.NewApp()
	cmd.Authors = []cli.Author{
		cli.Author{
			Name:  "JulzDiverse",
			Email: "julian.skupnjak@gmail.com",
		},
	}
	cmd.Name = "Aviator"
	cmd.Usage = "Navigate to a aviator.yml file and run aviator"
	cmd.Version = "0.20.0"
	cmd.Flags = getFlags()
	return cmd
}

func getFlags() []cli.Flag {
	var flags []cli.Flag
	flags = []cli.Flag{
		cli.StringFlag{
			Name:  "file, f",
			Value: "aviator.yml",
			Usage: "Specifies a path to an aviator yaml",
		},
		cli.BoolFlag{
			Name:  "verbose, vv",
			Usage: "prints warnings",
		},
		cli.BoolFlag{
			Name:  "silent, s",
			Usage: "silent mode (no prints)",
		},
		cli.StringSliceFlag{
			Name:  "var",
			Usage: "provides a variable to an aviator file: [key=value]",
		},
		cli.BoolFlag{
			Name:  "curly-braces, b",
			Usage: "allow {{}} syntax in yaml files",
		},
	}
	return flags
}
