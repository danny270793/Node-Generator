package git

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/danny270793/nodegenerator/npm"
	"github.com/danny270793/nodegenerator/shell"
)

func Init(project npm.Project) error {
	log.Printf("git.Init\n")
	_, err := shell.Execute(fmt.Sprintf("cd %s && git init", project.Path))
	if err != nil {
		return err
	}
	return nil
}

func GenerateGitignore(project npm.Project, ignores []string) error {
	log.Printf("git.GenerateGitignore\n")
	err := os.WriteFile(fmt.Sprintf("%s/.gitignore", project.Path), []byte(strings.Join(ignores, "\n")), 0774)
	if err != nil {
		return err
	}
	return nil
}

func Version() (string, error) {
	log.Printf("git.Version\n")
	stdout, err := shell.Execute("git --version")
	if err != nil {
		return "", err
	}
	return strings.Trim(strings.ReplaceAll(stdout, "git version ", ""), ""), nil
}
