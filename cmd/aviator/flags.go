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
	cmd.Version = "0.11.0"
	cmd.Flags = getFlags()
	return cmd
}

func getFlags() []cli.Flag {
	var flags []cli.Flag
	flags = []cli.Flag{
		cli.StringFlag{
			Name:  "file, f",
			Value: "aviator.yml",
			Usage: "Path to a AVIATOR YAML",
		},
		cli.BoolFlag{
			Name:  "verbose, vv",
			Usage: "print warnings",
		},
		cli.BoolFlag{
			Name:  "silent, s",
			Usage: "silent mode (no prints)",
		},
	}
	return flags
}
