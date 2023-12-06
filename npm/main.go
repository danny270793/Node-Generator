package npm

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/danny270793/nodegen/shell"
)

func regenerateFolder(path string) error {
	log.Printf("npm.regenerateFolder\n")
	command := fmt.Sprintf("rm -rf %s && mkdir -p %s", path, path)
	_, err := shell.Execute(command)
	if err != nil {
		return err
	}
	return nil
}

type PackageJson struct {
	Author          string            `json:"author"`
	Description     string            `json:"description"`
	DevDependencies map[string]string `json:"devDependencies"`
	Dependencies    map[string]string `json:"dependencies"`
	Keywords        []string          `json:"keywords"`
	License         string            `json:"license"`
	Main            string            `json:"main"`
	Type            string            `json:"type"`
	Name            string            `json:"name"`
	Scripts         map[string]string `json:"scripts"`
	Version         string            `json:"version"`
}

type Project struct {
	Path        string
	PackageJson PackageJson
}

func NewProject(path string) (Project, error) {
	log.Printf("npm.NewProject\n")
	err := regenerateFolder(path)
	if err != nil {
		return Project{}, nil
	}

	project := Project{
		Path: path,
	}

	err = project.Init()
	if err != nil {
		return Project{}, nil
	}

	err = project.AddScript("test", "node --test ./build")
	if err != nil {
		return Project{}, err
	}

	return project, nil
}

func (project *Project) Init() error {
	log.Printf("npm.Init\n")
	command := fmt.Sprintf("cd %s && npm init -y", project.Path)
	_, err := shell.Execute(command)
	if err != nil {
		return err
	}

	packageJson, err := project.LoadPackageJson()
	if err != nil {
		return err
	}
	project.PackageJson = packageJson
	//project.PackageJson.Type = "module"
	project.PackageJson.Main = "./build/src/index.js"
	err = project.WritePackageJson()
	if err != nil {
		return err
	}

	return nil
}

func (projecr *Project) LoadPackageJson() (PackageJson, error) {
	log.Printf("npm.LoadPackageJson\n")
	packageJsonPath := fmt.Sprintf("%s/package.json", projecr.Path)
	packageJsonContent, err := os.ReadFile(packageJsonPath)
	if err != nil {
		return PackageJson{}, nil
	}

	var packageJson PackageJson
	err = json.Unmarshal(packageJsonContent, &packageJson)
	if err != nil {
		return PackageJson{}, err
	}

	if packageJson.Dependencies == nil {
		packageJson.Dependencies = make(map[string]string)
	}

	return packageJson, nil
}

func (project *Project) WritePackageJson() error {
	log.Printf("npm.WritePackageJson\n")
	text, err := json.MarshalIndent(project.PackageJson, "", "  ")
	if err != nil {
		return err
	}

	packageJsonPath := fmt.Sprintf("%s/package.json", project.Path)
	err = os.WriteFile(packageJsonPath, text, 0775)
	if err != nil {
		return err
	}

	return nil
}

func (project *Project) Install(dependency string, isDev bool) error {
	log.Printf("npm.Install\n")
	mode := ""
	if isDev {
		mode = "--save-dev"
	} else {
		mode = "--save"
	}
	command := fmt.Sprintf("cd %s && npm install %s %s", project.Path, mode, dependency)
	_, err := shell.Execute(command)
	if err != nil {
		return err
	}

	packageJson, err := project.LoadPackageJson()
	if err != nil {
		return nil
	}
	project.PackageJson = packageJson

	return nil
}

func (project *Project) AddScript(name string, action string) error {
	log.Printf("npm.AddScript\n")
	project.PackageJson.Scripts[name] = action
	packageJsonString, err := json.MarshalIndent(project.PackageJson, "", "  ")
	if err != nil {
		return err
	}
	packageJsonPath := fmt.Sprintf("%s/package.json", project.Path)
	os.WriteFile(packageJsonPath, []byte(packageJsonString), 0774)
	return nil
}

func (project *Project) GenerateDotEnv(environment map[string]string) error {
	log.Printf("npm.GenerateDotEnv\n")
	keys := make([]string, 0, len(environment))
	for key := range environment {
		keys = append(keys, fmt.Sprintf("%s=%s", key, environment[key]))
	}

	err := os.WriteFile(fmt.Sprintf("%s/.env", project.Path), []byte(strings.Join(keys, "\n")), 0774)
	if err != nil {
		return err
	}

	err = os.WriteFile(fmt.Sprintf("%s/.env.sample", project.Path), []byte(strings.Join(keys, "\n")), 0774)
	if err != nil {
		return err
	}

	return nil
}

func NodeVersion() (string, error) {
	log.Printf("npm.NodeVersion\n")
	stdout, err := shell.Execute("node --version")
	if err != nil {
		return "", err
	}
	return strings.Trim(stdout, ""), nil
}

func Version() (string, error) {
	log.Printf("npm.Version\n")
	stdout, err := shell.Execute("npm --version")
	if err != nil {
		return "", err
	}
	return strings.Trim(stdout, ""), nil
}
