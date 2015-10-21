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

func DeleteRemoteBranch(branch string, prompt bool) {
	var err error

	if prompt == true {
		var input string
		fmt.Printf("remove branch %v [y/N]: ", branch)
		fmt.Scanln(&input)

		if input == "y" || input == "Y" {
			_, err = exec.Command("git", "push", "origin", fmt.Sprintf(":%v", branch)).Output()
		}

	} else {
		_, err = exec.Command("git", "push", "origin", fmt.Sprintf(":%v", branch)).Output()
	}

	if err != nil {
		log.Fatalf("Error deleting branch %v.\n", branch)
	}
}

func DeleteLocalBranch(branch string, prompt bool) {
	var err error

	if prompt == true {
		var input string
		fmt.Printf("remove local branch %v [y/N]: ", branch)
		fmt.Scanln(&input)

		if input == "y" || input == "Y" {
			_, err = exec.Command("git", "branch", "-D", branch).Output()
		}

	} else {
		_, err = exec.Command("git", "branch", "-D", branch).Output()
	}

	if err != nil {
		log.Println("Error '%v' does not exist.\n", branch)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "grb"
	app.Version = "0.1.2"
	app.Usage = "grb [global options] \"search terms\""

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "delete, d",
			Usage: "Deletes remote branches return by search result.",
		},
		cli.BoolFlag{
			Name:  "local, l",
			Usage: "Used with -d or --delete to remove local branch along with remote branches.",
		},
		cli.BoolFlag{
			Name:  "no-prompt, f",
			Usage: "Attempt to delete without confirmation.",
		},
		cli.StringFlag{
			Name:  "remote, r",
			Usage: "Alters git remote searched by grb. defaults to origin if flag is not provided.",
		},
	}

	remote := Remote{
		Name: "origin",
	}

	app.Action = func(c *cli.Context) {
		var searchString string
		var promptFlag = true

		if c.String("no-prompt") == "true" {
			promptFlag = false
			fmt.Println(promptFlag)
		}

		if len(c.String("remote")) > 0 {
			remote.Name = c.String("remote")
		}

		branches := remote.Fetch()

		if len(c.Args()) > 0 {
			searchString = c.Args()[0]
		}

		if c.String("delete") == "true" {
			if len(c.Args()) > 0 && len(c.Args()[0]) > 0 {
				for _, b := range Filter(branches, searchString) {
					DeleteRemoteBranch(b, promptFlag)
				}
				if c.String("local") == "true" {
					local := new(Local)
					localBranches := local.Fetch()
					for _, b := range Filter(localBranches, searchString) {
						DeleteLocalBranch(b, promptFlag)
					}
				}
			} else {
				log.Println("Please provide search terms.")
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
