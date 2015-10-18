package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"os/exec"
	"strings"
)

type branch struct {
	Name string
	Repo string
}

var heads []string

func parseHeads(remote string) {
	heads = nil
	out, err := exec.Command("git", "ls-remote", "--heads", remote).Output()
	if err != nil {
		panic(err)
	}

	raw := strings.Fields(string(out))
	for _, field := range raw {
		if strings.Contains(field, "refs/heads") {
			heads = append(heads, field)
		}
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "grb"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "grb list",
			Action: func(c *cli.Context) {
				parseHeads("origin")
				for _, h := range heads {
					fmt.Println(h)
				}
			},
		},
	}
	app.Run(os.Args)
}
