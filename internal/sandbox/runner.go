package sandbox

import "os/exec"

type Runner interface {
	Command(
		command []string,
		workspace string,
		wallTimeS int,
	) *exec.Cmd
}
