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

	var cleanArgs []string
	var outFile string
	var appendMode bool

	// 1. Analisa os argumentos para intercetar redirecionamentos
	for i := 0; i < len(args); i++ {
		if args[i] == ">" {
			if i+1 < len(args) {
				outFile = args[i+1]
				i++ // Pula o nome do arquivo para não o passar ao comando
			} else {
				return fmt.Errorf("erro de sintaxe: esperado nome do arquivo após '>'")
			}
			appendMode = false
		} else if args[i] == ">>" {
			if i+1 < len(args) {
				outFile = args[i+1]
				i++
			} else {
				return fmt.Errorf("erro de sintaxe: esperado nome do arquivo após '>>'")
			}
			appendMode = true
		} else {
			// Argumento normal
			cleanArgs = append(cleanArgs, args[i])
		}
	}

	if len(cleanArgs) == 0 {
		return nil // Caso digitem apenas "> arquivo.txt"
	}

	// 2. Prepara o comando apenas com os argumentos limpos
	cmd := exec.Command(cleanArgs[0], cleanArgs[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr // Deixamos o stderr na tela por enquanto

	// 3. Configura o Stdout (Tela ou Arquivo)
	if outFile != "" {
		flags := os.O_WRONLY | os.O_CREATE
		if appendMode {
			flags |= os.O_APPEND
		} else {
			flags |= os.O_TRUNC
		}

		// Abre (ou cria) o arquivo a nível do SO
		file, err := os.OpenFile(outFile, flags, 0644)
		if err != nil {
			return fmt.Errorf("erro ao abrir arquivo de redirecionamento: %v", err)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {

			}
		}(file)

		// O comando vai escrever neste arquivo em vez do os.Stdout
		cmd.Stdout = file
	} else {
		cmd.Stdout = os.Stdout
	}

	currentCmd = cmd

	// 4. Executa o comando
	if err := cmd.Start(); err != nil {
		currentCmd = nil
		return err
	}

	if err := cmd.Wait(); err != nil {
		currentCmd = nil

		// Ignora erros de sinal SIGINT (Ctrl+C)
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
		_ = syscall.Kill(currentCmd.Process.Pid, syscall.SIGINT)
		currentCmd = nil
	}
}
