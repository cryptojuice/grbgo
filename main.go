package main

import (
	"fmt"
	"log"
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
		log.Fatalf("Error fetching from '%v' remote.\n", r.Name)
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
		log.Fatalf("Error deleting branch %v.\n", branch)
	}
}

func DeleteLocalBranch(prompt bool, branch string) {
	_, err := exec.Command("git", "branch", "-D", branch).Output()
	if err != nil {
		log.Fatalf("Error deleting local branch %v.\n", branch)
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
		cli.StringFlag{
			Name: "remote, r",
		},
	}

	remote := Remote{
		Name: "origin",
	}

	app.Action = func(c *cli.Context) {
		var searchString string
		branches := remote.Fetch()

		if len(c.Args()) > 0 {
			searchString = c.Args()[0]
		}

		if len(c.String("remote")) > 0 {
			remote.Name = c.String("remote")
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
			Usage:   "grb list, grb list 'possible-branch-name', grb -r 'remoteName' list",
			Action: func(c *cli.Context) {
				if len(c.GlobalString("remote")) > 0 {
					remote.Name = c.GlobalString("remote")
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
