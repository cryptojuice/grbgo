package repositories

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type Local struct {
	Branches []string
}

func (l *Local) Fetch() []string {
	out, err := exec.Command("git", "branch").Output()
	if err != nil {
		log.Fatalf("Error fetching local branches.\n")
	}

	l.Branches = strings.Fields(string(out))
	return l.Branches
}

func (r *Local) DeleteBranch(branch string, promptForDeletion bool) {
	var err error

	if promptForDeletion == true {
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
