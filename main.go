package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/codegangsta/cli"
)

type Remote struct {
	Name     string
	Branches []string
}

func (r *Remote) Fetch() []string {
	out, err := exec.Command("git", "ls-remote", "--heads", r.Name).Output()
	if err != nil {
		panic(err)
	}

	raw := strings.Fields(string(out))
	for _, field := range raw {
		if strings.Contains(field, "refs/heads") {
			r.Branches = append(r.Branches, field)
		}
	}

	return r.Branches
}

func Filter(branches []string, searchString string) []string {
	filtered := branches[:0]
	for _, b := range branches {
		if strings.Contains(b, searchString) {
			filtered = append(filtered, b)
		}
	}
	return filtered
}

func DeleteRemoteBranch(prompt bool, branch string) {
	_, err := exec.Command("git", "push", "origin", fmt.Sprintf(":%v", branch)).Output()
	if err != nil {
		panic(err)
	}
}

func DeleteLocalBranch(prompt bool, branch string) {
	_, err := exec.Command("git", "branch", "-D", branch).Output()
	if err != nil {
		panic(err)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "grb"
	app.Version = "0.1.0"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name: "delete, d",
		},
		cli.BoolFlag{
			Name: "local, l",
		},
	}

	app.Action = func(c *cli.Context) {
		var searchString string
		remote := Remote{
			Name: "origin",
		}
		branches := remote.Fetch()

		if len(c.Args()) > 0 {
			searchString = c.Args()[0]
		}

		if c.String("local") == "true" {
			for _, b := range Filter(branches, searchString) {
				DeleteLocalBranch(false, b)
			}
		}
		if c.String("delete") == "true" {
			if len(c.Args()) > 0 {
				for _, b := range Filter(branches, searchString) {
					DeleteRemoteBranch(false, b)
				}
			}
		}
	}

	app.Commands = []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "grb list",
			Action: func(c *cli.Context) {
				remote := Remote{
					Name: "origin",
				}
				branches := remote.Fetch()
				if len(c.Args()) > 0 {
					searchString := c.Args()[0]
					for _, b := range Filter(branches, searchString) {
						fmt.Println(b)
					}
				} else {
					for _, b := range branches {
						fmt.Println(b)
					}
				}
			},
		},
	}
	app.Run(os.Args)
}
