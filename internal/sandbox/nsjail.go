package sandbox

import (
	"os/exec"
	"strconv"
)

type NSJailRunner struct{}

func (r NSJailRunner) Command(
	command []string,
	workspace string,
	wallTimeS int,
) *exec.Cmd {

	args := []string{
		"--mode",
		"o",

		"--chroot",
		"/",

		"--cwd",
		workspace,

		"--time_limit",
		strconv.Itoa(wallTimeS),

		"--",
	}

	args = append(args, command...)

	return exec.Command(
		"nsjail",
		args...,
	)
}
