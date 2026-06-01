package sandbox

import "os/exec"

type NativeRunner struct{}

func (r NativeRunner) Command(
	command []string,
	workspace string,
	wallTimeS int,
) *exec.Cmd {

	return exec.Command(
		command[0],
		command[1:]...,
	)
}
