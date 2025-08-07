package main

import (
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/FelipeFelipeRenan/gosh/internal/builtin"
	"github.com/FelipeFelipeRenan/gosh/internal/executor"
	"github.com/FelipeFelipeRenan/gosh/internal/history"
	"github.com/FelipeFelipeRenan/gosh/internal/parser"
	"github.com/FelipeFelipeRenan/gosh/internal/signals"

	"golang.org/x/term"
)
func main() {
	signals.SetupSignalHandlers()

	usr, err := user.Current()
	if err != nil {
		fmt.Println("erro de usuário: ", err)
	}

	// Ativa raw apenas uma vez
	initialState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("erro ao configurar terminal: ", err)
		return
	}
	defer term.Restore(int(os.Stdin.Fd()), initialState)

	history := history.New()

	for {
		fmt.Printf("\r\ngosh | %s > ", usr.Username)

		var input []rune
		history.ResetPos()

	readLoop:
		for {
			rb := make([]byte, 1)
			_, err := os.Stdin.Read(rb)
			if err != nil {
				fmt.Println("(gosh) erro ao ler entrada:", err)
				return
			}
			b := rb[0]

			switch b {
			case 3: // Ctrl+C
				fmt.Println("^C")
				input = nil
				break readLoop
			case 13: // Enter
				fmt.Println()
				break readLoop
			case 127: // Backspace
				if len(input) > 0 {
					input = input[:len(input)-1]
					fmt.Print("\b \b")
				}
			case 27: // Escape sequence
				seq := make([]byte, 2)
				os.Stdin.Read(seq)
				if seq[0] == '[' {
					switch seq[1] {
					case 'A': // Up arrow
						prev := history.Prev()
						clearLine(len(input))
						input = []rune(prev)
						fmt.Print(string(input))
					case 'B': // Down arrow
						next := history.Next()
						clearLine(len(input))
						input = []rune(next)
						fmt.Print(string(input))
					}
				}
			default:
				fmt.Print(string(b))
				input = append(input, rune(b))
			}
		}

		cmd := strings.TrimSpace(string(input))
		if cmd == "" {
			continue
		}
		history.Add(cmd)

		args := parser.Parse(cmd)

		handled, err := builtin.Exec(args)
		if err != nil {
			fmt.Println("(gosh) erro no comando interno:", err)
			continue
		}
		if handled {
			continue
		}

		// ⚠️ Restaura terminal antes de comando externo
		term.Restore(int(os.Stdin.Fd()), initialState)

		if err := executor.Exec(args); err != nil {
			fmt.Println("(gosh) erro ao executar comando:", err)
		}

		// ⚠️ Reativa modo raw após execução
		initialState, err = term.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			fmt.Println("erro ao reativar modo raw:", err)
			return
		}
	}
}


func clearLine(n int) {
	for i := 0; i < n; i++ {
		fmt.Print("\b \b")
	}
}
