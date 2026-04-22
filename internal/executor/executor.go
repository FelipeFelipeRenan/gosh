package executor

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

var currentCmd *exec.Cmd

func Exec(args []string) error {
	if len(args) == 0 {
		return nil
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	currentCmd = cmd

	if err := cmd.Start(); err != nil {
		currentCmd = nil
		return err
	}

	if err := cmd.Wait(); err != nil {
		currentCmd = nil

		// Ignora erros de sinal SIGINT
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				if status.Signaled() && status.Signal() == syscall.SIGINT {
					fmt.Println()
					return nil
				}
			}
		}
		return err
	}

	currentCmd = nil
	return nil
}

func InterruptCurrentCommand() {
	if currentCmd != nil && currentCmd.Process != nil {
		// Envia SIGINT para todo o grupo de processos
		_ = syscall.Kill(currentCmd.Process.Pid, syscall.SIGINT)
		currentCmd = nil
	}
}
