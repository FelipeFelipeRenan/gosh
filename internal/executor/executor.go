package executor

import (
	"os"
	"os/exec"
)


func Exec(args []string) error{
	if len(args) == 0{
		return nil
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}