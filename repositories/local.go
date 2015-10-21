package repositories

import (
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
