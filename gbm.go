package main

import (
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "Git Branch Manager"
	app.Version = "0.0.1"

	// app.Flags = []cli.Flag{
	// 	cli.StringFlag{
	// 		Name: "list",
	// 	},
	// }

	app.Commands = []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "list remote branches",
			Action: func(c *cli.Context) {
				println("listing remote branches")
			},
		},
	}

	app.Action = func(c *cli.Context) {
	}

	app.Run(os.Args)
}
