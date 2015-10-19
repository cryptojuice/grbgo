package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/codegangsta/cli"
)

type Head struct {
	Branches []string
}

func (h *Head) PopulateBranches(remote string) {
	out, err := exec.Command("git", "ls-remote", "--heads", remote).Output()
	if err != nil {
		panic(err)
	}

	raw := strings.Fields(string(out))
	for _, field := range raw {
		if strings.Contains(field, "refs/heads") {
			h.Branches = append(h.Branches, field)
		}
	}
}

func (h *Head) Filter(searchString string) []string {
	fb := h.Branches[:0]
	for _, b := range h.Branches {
		if strings.Contains(b, searchString) {
			fb = append(fb, b)
		}
	}
	return fb
}

func Delete(prompt bool) {
}

func main() {
	app := cli.NewApp()
	app.Name = "grb"
	app.Version = "0.1.0"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name: "delete, d",
		},
	}

	app.Action = func(c *cli.Context) {
		h := Head{}
		h.PopulateBranches("origin")
		if len(c.String("delete")) > 0 {
			searchString := c.Args()[0]
			for _, b := range h.Filter(searchString) {
				_, err := exec.Command("git", "push", "origin", fmt.Sprintf(":%v", b)).Output()
				if err != nil {
					panic(err)
				}
				fmt.Println(b)
			}
		}
	}

	app.Commands = []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "grb list",
			Action: func(c *cli.Context) {
				h := Head{}
				h.PopulateBranches("origin")
				if len(c.Args()) > 0 {
					searchString := c.Args()[0]
					for _, b := range h.Filter(searchString) {
						fmt.Println(b)
					}
				} else {
					for _, b := range h.Branches {
						fmt.Println(b)
					}
				}
			},
		},
	}
	app.Run(os.Args)
}
