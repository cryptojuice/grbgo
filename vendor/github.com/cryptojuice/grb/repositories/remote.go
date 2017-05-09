package repositories

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
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

func (r *Remote) DeleteBranch(branch string, promptForDeletion bool) {
	var err error

	if promptForDeletion == true {
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
