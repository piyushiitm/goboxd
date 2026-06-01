package sandbox

import "os/exec"

type NSJailRunner struct{}

func (r NSJailRunner) Command(
	command []string,
	workspace string,
	wallTimeS int,
) *exec.Cmd {

	args := []string{
		"--mode",
		"o",
		"--cwd",
		workspace,
		"--time_limit",
		string(rune(wallTimeS)),
		"--",
	}

	args = append(args, command...)

	return exec.Command(
		"nsjail",
		args...,
	)
}
