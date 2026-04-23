package main

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/FelipeFelipeRenan/gosh/internal/builtin"
	"github.com/FelipeFelipeRenan/gosh/internal/executor"
	"github.com/FelipeFelipeRenan/gosh/internal/history"
	"github.com/FelipeFelipeRenan/gosh/internal/parser"
	"github.com/FelipeFelipeRenan/gosh/internal/signals"
	"github.com/FelipeFelipeRenan/gosh/internal/trie"

	"golang.org/x/term"
)

func main() {
	signals.SetupSignalHandlers()

	usr, err := user.Current()
	if err != nil {
		fmt.Println("erro de usuário: ", err)
	}

	// Salva o estado original do terminal APENAS UMA VEZ
	fd := int(os.Stdin.Fd())
	initialState, err := term.GetState(fd)
	if err != nil {
		fmt.Println("erro ao obter estado do terminal: ", err)
		return
	}
	// Garante a restauração ao sair do programa
	defer func(fd int, oldState *term.State) {
		err := term.Restore(fd, oldState)
		if err != nil {
			fmt.Println("erro ao restore terminal state: ", err)
		}
	}(fd, initialState)

	historyFile := ""
	if usr != nil {
		historyFile = filepath.Join(usr.HomeDir, ".gosh_history")
	}
	historyInstace := history.New(historyFile)
	cmdTrie := trie.New()
	loadBinariesIntoTrie(cmdTrie)

	for {
		pwd, _ := os.Getwd()

		// 1. Garante que o terminal está normal para o prompt
		err := term.Restore(fd, initialState)
		if err != nil {
			return
		}
		fmt.Printf("gosh | %s at %s > ", usr.Username, pwd)

		// 2. Ativa modo Raw para leitura de teclas
		rawState, err := term.MakeRaw(fd)
		if err != nil {
			fmt.Printf("\r\n(gosh) erro ao ativar modo raw: %v\n", err)
			return
		}

		var input []rune
		historyInstace.ResetPos()

	readLoop:
		for {
			rb := make([]byte, 1)
			_, err := os.Stdin.Read(rb)
			if err != nil {
				break readLoop
			}
			b := rb[0]

			switch b {
			case 3: // Ctrl+C
				input = nil
				fmt.Print("^C\r\n")
				break readLoop
			case 9:
				prefix := string(input)
				if prefix == "" {
					continue
				}
				suggestions := cmdTrie.SearchPrefix(prefix)
				if len(suggestions) == 1 {
					clearLine(len(input))
					input = []rune(suggestions[0])
					fmt.Print(string(input))
				} else if len(suggestions) > 1 {
					// Várias sugestões: imprime-as abaixo (como o bash)
					fmt.Printf("\r\n%s\r\n", strings.Join(suggestions, "  "))
					// Mostra o prompt e o que já foi digitado novamente
					fmt.Printf("gosh | %s at %s > %s", usr.Username, pwd, string(input))
				}
			case 13: // Enter (Carriage Return)
				fmt.Print("\r\n")
				break readLoop
			case 127, 8: // Backspace
				if len(input) > 0 {
					input = input[:len(input)-1]
					fmt.Print("\b \b") // Move volta, imprime espaço, move volta
				}
			case 27: // Escape sequences (Setas)
				seq := make([]byte, 2)
				_, err2 := os.Stdin.Read(seq)
				if err2 != nil {
					return
				}
				if seq[0] == '[' {
					switch seq[1] {
					case 'A': // Up arrow
						clearLine(len(input))
						prev := historyInstace.Prev()
						input = []rune(prev)
						fmt.Print(string(input))
					case 'B': // Down arrow
						clearLine(len(input))
						next := historyInstace.Next()
						input = []rune(next)
						fmt.Print(string(input))
					}
				}
			default:
				// CORREÇÃO AQUI: Use Printf para formatar o caractere
				fmt.Printf("%c", b)
				input = append(input, rune(b))
			}
		}

		// 3. Restaura o terminal IMEDIATAMENTE após a leitura
		err = term.Restore(fd, rawState)
		if err != nil {
			return
		}
		err = term.Restore(fd, initialState)
		if err != nil {
			return
		}

		cmd := strings.TrimSpace(string(input))
		if cmd == "" {
			continue
		}
		historyInstace.Add(cmd)

		args := parser.Parse(cmd)

		// Execução de comandos
		handled, err := builtin.Exec(args, historyInstace)
		if err != nil {
			if errors.Is(err, builtin.ErrExit) {
				return
			}
			fmt.Printf("(gosh) erro no comando interno: %v\n", err)
			continue
		}

		if !handled {
			if err := executor.Exec(args); err != nil {
				fmt.Printf("(gosh) erro ao executar comando: %v\n", err)
			}
		}
	}
}

func clearLine(n int) {
	// Apaga os caracteres da tela movendo o cursor para trás e sobrescrevendo com espaços
	for i := 0; i < n; i++ {
		fmt.Print("\b \b")
	}
}

func loadBinariesIntoTrie(t *trie.Trie) {
	pathVar := os.Getenv("PATH")
	dirs := filepath.SplitList(pathVar)

	for _, dir := range dirs {
		files, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, file := range files {
			if !file.IsDir() {
				t.Insert(file.Name())
			}
		}
	}
}
