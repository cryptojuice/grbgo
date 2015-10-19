package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/codegangsta/cli"
)

type Repository interface {
	Fetch() []string
}

type Remote struct {
	Name     string
	Branches []string
}

type Local struct {
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

func (l *Local) Fetch() []string {
	out, err := exec.Command("git", "branch").Output()
	if err != nil {
		log.Fatalf("Error fetching local branches.\n")
	}

	l.Branches = strings.Fields(string(out))
	return l.Branches
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
		log.Fatalf("Error '%v' may not exist locally.\n", branch)
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

		if len(c.String("remote")) > 0 {
			remote.Name = c.String("remote")
		}

		if c.String("local") == "true" {
			local := new(Local)
			localBranches := local.Fetch()
			for _, b := range Filter(localBranches, searchString) {
				DeleteLocalBranch(false, b)
			}
		}

		branches := remote.Fetch()

		if c.String("delete") == "true" {
			if len(c.Args()) > 0 {
				for _, b := range Filter(branches, searchString) {
					DeleteRemoteBranch(false, b)
				}
			}
		}

		if len(c.Args()) == 0 {
			for _, b := range branches {
				fmt.Println(string(b[11:]))
			}
		}

		if len(c.Args()) > 0 && c.String("local") == "false" && c.String("delete") == "false" {
			searchString = c.Args()[0]
			for _, b := range Filter(branches, searchString) {
				fmt.Println(string(b[11:]))
			}
		}

	}

	app.Run(os.Args)
}
