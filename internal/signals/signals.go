// signals/signals.go
package signals

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/FelipeFelipeRenan/gosh/internal/executor"
)

func SetupSignalHandlers() {
	signals := make(chan os.Signal, 1)

	// captura os sinais SIGINT(Ctrl+C) e SIGTSTP (Ctrl+Z)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTSTP) 

	go func() {
		for sig := range signals {
			switch sig {
			case syscall.SIGINT:
				// interrompe o comando em execução mas não fecha o shell
				executor.InterruptCurrentCommand()
				fmt.Print("\n(gosh) comando interrompido\n")
			case syscall.SIGTSTP:
				// ignora ctrl+z
				fmt.Print("\n(gosh) Ctrl+Z ignorado")
			}
		}	
	}()

}