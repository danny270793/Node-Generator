package typescript

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/danny270793/nodegenerator/npm"
	"github.com/danny270793/nodegenerator/shell"
)

func Configure(project *npm.Project) error {
	log.Printf("typescript.Configure\n")
	err := project.Install("typescript ts-node-dev @typescript-eslint/eslint-plugin", true)
	if err != nil {
		log.Printf("error instaling typescript\n")
		return err
	}

	err = project.AddScript("build", "tsc")
	if err != nil {
		return err
	}
	err = project.AddScript("start", "node ./build/src/index.js")
	if err != nil {
		return err
	}
	err = project.AddScript("start:watch", "ts-node-dev --respawn ./src/index.ts")
	if err != nil {
		return err
	}
	err = project.AddScript("test:watch", "ts-node-dev --respawn ./tests/index.test.ts")
	if err != nil {
		return err
	}

	_, err = shell.Execute(fmt.Sprintf("cd %s && npx tsc --init", project.Path))
	if err != nil {
		log.Printf("error initing typescript project\n")
		return err
	}

	tsconfigJson := fmt.Sprintf("%s/tsconfig.json", project.Path)
	bytes, err := os.ReadFile(tsconfigJson)
	if err != nil {
		return err
	}
	newFile := ""
	for _, line := range strings.Split(string(bytes), "\n") {
		if strings.Contains(line, "// \"outDir\": \"./\"") {
			newLine := strings.ReplaceAll(line, "// \"outDir\": \"./\"", "\"outDir\": \"./build\"")
			newFile += newLine
		} else if strings.Contains(line, "// \"declaration\": ") {
			newLine := "    \"declaration\": true,                              /* Generate .d.ts files from TypeScript and JavaScript files in your project. */"
			newFile += newLine
		} else {
			newFile += line
		}
		newFile += "\n"
	}
	err = os.WriteFile(tsconfigJson, []byte(newFile), 0775)
	if err != nil {
		return err
	}

	return nil
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GenerateCode(project *npm.Project) error {
	log.Printf("typescript.GenerateCode\n")

	//main
	srcPath := fmt.Sprintf("%s/src", project.Path)
	exists, err := pathExists(srcPath)
	if err != nil {
		return err
	}
	if !exists {
		os.RemoveAll(srcPath)
	}
	os.Mkdir(srcPath, 0775)

	mainCode := "async function main(): Promise<void> {\n  console.log('Hello world')\n}\n\nmain().catch(console.error)\n"
	indexPath := fmt.Sprintf("%s/index.ts", srcPath)

	mainFile, err := os.OpenFile(indexPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0775)
	if err != nil {
		return err
	}
	defer mainFile.Close()
	_, err = mainFile.WriteString(mainCode)
	if err != nil {
		return err
	}

	//test
	testsPath := fmt.Sprintf("%s/tests", project.Path)
	testExists, err := pathExists(testsPath)
	if err != nil {
		return err
	}
	if !testExists {
		os.RemoveAll(testsPath)
	}
	os.Mkdir(testsPath, 0775)

	testCode := "import { describe, it } from 'node:test'\nimport assert from 'node:assert'\n\ndescribe('module', () => {\n  it('should be equals', () => {\n    assert.equal(1, 1)\n  })\n  it('should not be equals', () => {\n    assert.notEqual(1, 0)\n  })\n})\n"
	testPath := fmt.Sprintf("%s/index.test.ts", testsPath)

	testFile, err := os.OpenFile(testPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0775)
	if err != nil {
		return err
	}
	defer testFile.Close()
	_, err = testFile.WriteString(testCode)
	if err != nil {
		return err
	}

	return nil
}
