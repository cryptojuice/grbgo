package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cryptojuice/grb/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/cryptojuice/grb/repositories"
)

func Filter(branches []string, searchString string) []string {
	filtered := branches[:0]
	for _, b := range branches {
		if strings.Contains(b, searchString) {
			filtered = append(filtered, b)
		}
	}
	return filtered
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

	app.Action = func(c *cli.Context) {
		var promptFlag = true
		var deleteRemoteFlag = false
		var deleteLocalFlag = false
		var searchString string

		var remoteRepository = repositories.Remote{
			Name: "origin",
		}
		var localRepository = repositories.Local{}

		if c.String("no-prompt") == "true" {
			promptFlag = false
		}

		if len(c.String("remote")) > 0 {
			remoteRepository.Name = c.String("remote")
		}

		if c.String("delete") == "true" {
			deleteRemoteFlag = true
		}

		if c.String("local") == "true" {
			deleteLocalFlag = true
		}

		branches := remoteRepository.Fetch()
		localBranches := localRepository.Fetch()

		if len(c.Args()) > 0 {
			searchString = c.Args()[0]
		}

		if deleteRemoteFlag == true {
			if len(c.Args()) > 0 && len(c.Args()[0]) > 0 {
				for _, b := range Filter(branches, searchString) {
					remoteRepository.DeleteBranch(b, promptFlag)
				}

				if deleteLocalFlag == true {
					for _, b := range Filter(localBranches, searchString) {
						localRepository.DeleteBranch(b, promptFlag)
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

		if len(c.Args()) > 0 && deleteLocalFlag == false && deleteRemoteFlag == false {
			searchString = c.Args()[0]
			for _, b := range Filter(branches, searchString) {
				fmt.Println(string(b[11:]))
			}
		}
	}

	app.Run(os.Args)
}
