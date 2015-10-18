package main

import (
	"fmt"
	_ "github.com/codegangsta/cli"
	"os/exec"
	"strings"
)

type branch struct {
	Name string
	Repo string
}

var heads []string

func parseHeads(remote string) {
	out, err := exec.Command("git", "ls-remote", "--heads", remote).Output()
	if err != nil {
		panic(err)
	}

	raw := strings.Fields(string(out))
	for _, field := range raw {
		if strings.Contains(field, "refs/heads") {
			fmt.Println(field)
		}
	}
}

func main() {
	parseHeads("origin")
	// app := cli.NewApp()
	// app.Name = "Git Branch Manager"
	// app.Version = "0.0.1"
	// app.Commands = []cli.Command{
	// 	{
	// 		Name:    "list",
	// 		Aliases: []string{"l"},
	// 		Usage:   "grb list",
	// 		Action: func(c *cli.Context) {
	// 		},
	// 	},
	// }
	// app.Run(os.Args)
}
