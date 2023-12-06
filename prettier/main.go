package prettier

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/danny270793/nodegenerator/npm"
)

type Prettierrc struct {
	Semi        bool `json:"semi"`
	SingleQuote bool `json:"singleQuote"`
	TabWidth    int  `json:"tabWidth"`
}

func Configure(project *npm.Project, prettierrc Prettierrc, ignores []string) error {
	log.Printf("prettier.Configure\n")
	err := project.Install("prettier", true)
	if err != nil {
		return err
	}
	err = project.AddScript("format", "prettier . --write")
	if err != nil {
		return err
	}

	text, err := json.MarshalIndent(prettierrc, "", "  ")
	if err != nil {
		log.Printf("error marshaling config\n")
		return err
	}
	err = os.WriteFile(fmt.Sprintf("%s/.prettierrc.json", project.Path), []byte(text), 0774)
	if err != nil {
		log.Printf("error writting prettierrc\n")
		return err
	}
	err = os.WriteFile(fmt.Sprintf("%s/.prettierignore", project.Path), []byte(strings.Join(ignores, "\n")), 0774)
	if err != nil {
		log.Printf("error writting prettierignore\n")
		return err
	}
	return nil
}
