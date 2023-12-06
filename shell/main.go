package shell

import (
	"bytes"
	"errors"
	"os/exec"
)

func Execute(command string) (string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	shell := exec.Command("bash", "-c", command)
	shell.Stdout = &stdout
	shell.Stderr = &stderr
	err := shell.Run()
	if err != nil {
		return "", errors.New(stderr.String())
	}
	return stdout.String(), nil
}
