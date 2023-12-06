package main

import (
	"fmt"
	"os"

	"github.com/danny270793/nodegenerator/eslint"
	"github.com/danny270793/nodegenerator/git"
	"github.com/danny270793/nodegenerator/npm"
	"github.com/danny270793/nodegenerator/prettier"
	"github.com/danny270793/nodegenerator/typescript"
)

func CheckVersions() error {
	node, nodeErr := npm.NodeVersion()
	npm, npmErr := npm.Version()
	git, gitErr := git.Version()

	fmt.Printf("node: %s\nnpm: %s\ngit: %s\n", node, npm, git)

	if nodeErr != nil {
		return npmErr
	}
	if npmErr != nil {
		return npmErr
	}
	if gitErr != nil {
		return gitErr
	}

	return nil
}

var VERSION string = "1.0.6"

func main() {
	firstArgument := os.Args[1]
	if firstArgument == "--help" {
		fmt.Printf("example:\n\nodegenerator create ./path/to/node/project\n")
		os.Exit(0)
	} else if firstArgument == "version" {
		fmt.Println(VERSION)
		os.Exit(0)
	} else if firstArgument == "create" {
		path := os.Args[2]

		err := CheckVersions()
		if err != nil {
			panic(err)
		}

		//node
		project, err := npm.NewProject(path)
		if err != nil {
			panic(err)
		}
		environment := map[string]string{
			"NODE_ENV": "development",
		}
		err = project.GenerateDotEnv(environment)
		if err != nil {
			panic(err)
		}

		//git
		err = git.Init(project)
		if err != nil {
			panic(err)
		}
		ignores := []string{
			"build",
			"node_modules",
			".env",
		}
		err = git.GenerateGitignore(project, ignores)
		if err != nil {
			panic(err)
		}

		//eslint
		eslintConfiguration := eslint.Eslintrc{
			Root:    true,
			Parser:  "@typescript-eslint/parser",
			Plugins: []string{"@typescript-eslint"},
			Extends: []string{
				"eslint:recommended",
				"plugin:@typescript-eslint/eslint-recommended",
				"plugin:@typescript-eslint/recommended",
			},
			Rules: map[string]int{
				"@typescript-eslint/no-inferrable-types": 0,
				"@typescript-eslint/no-explicit-any":     0,
			},
		}
		eslintIgnores := []string{"node_modules", "build"}
		err = eslint.Configure(&project, eslintConfiguration, eslintIgnores)
		if err != nil {
			panic(err)
		}

		//prettier
		prettierConfiguration := prettier.Prettierrc{
			Semi:        false,
			SingleQuote: true,
			TabWidth:    2,
		}
		prettierIgnores := []string{"node_modules", "build"}
		err = prettier.Configure(&project, prettierConfiguration, prettierIgnores)
		if err != nil {
			panic(err)
		}

		//typescript
		err = typescript.Configure(&project)
		if err != nil {
			panic(err)
		}
		err = typescript.GenerateCode(&project)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Printf("invalid action %s\n\n", firstArgument)
		fmt.Printf("example:\n\nodegenerator ./path/to/node/project\n")
		os.Exit(1)
	}
}
