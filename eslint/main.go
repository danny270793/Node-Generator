package eslint

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/danny270793/nodegen/npm"
)

type Eslintrc struct {
	Root    bool           `json:"root"`
	Parser  string         `json:"parser"`
	Plugins []string       `json:"plugins"`
	Extends []string       `json:"extends"`
	Rules   map[string]int `json:"rules"`
}

func Configure(project *npm.Project, eslintrc Eslintrc, ignores []string) error {
	log.Printf("eslint.Configure\n")
	err := project.Install("eslint", true)
	if err != nil {
		return err
	}
	err = project.AddScript("lint", "eslint . --ext .ts")
	if err != nil {
		return err
	}

	text, err := json.MarshalIndent(eslintrc, "", "  ")
	if err != nil {
		log.Printf("error marshaling config\n")
		return err
	}
	err = os.WriteFile(fmt.Sprintf("%s/.eslintrc.json", project.Path), []byte(text), 0774)
	if err != nil {
		log.Printf("error writting eslintrc\n")
		return err
	}
	err = os.WriteFile(fmt.Sprintf("%s/.eslintignore", project.Path), []byte(strings.Join(ignores, "\n")), 0774)
	if err != nil {
		log.Printf("error writting eslintignore\n")
		return err
	}
	return nil
}
