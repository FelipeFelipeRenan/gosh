// executor/executor.go
package executor

import (
	"errors"
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
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	currentCmd = cmd // guarda o cmd atual
	
	err := cmd.Run()
	currentCmd = nil // limpa após a execução

	if err != nil {
		// ignora erros se forem causados por sinal de interrupção
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			status := exitErr.Sys().(syscall.WaitStatus)
			if status.Signaled() && status.Signal() == syscall.SIGINT {
				// comando interrompido com Ctrl+C
				return nil
			}
		}
	}
	return err
}

// Exposed para o signal handler
func InterruptCurrentCommand()  {
	if currentCmd != nil && currentCmd.Process != nil{
		// envia o sinal SIGINT para o grupo de processos do comando
		_ = syscall.Kill(-currentCmd.Process.Pid, syscall.SIGINT)
	}
}
