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
	go loadBinariesIntoTrie(cmdTrie)

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
				line := string(input)
				parts := strings.Fields(line)

				isCommand := len(parts) == 0 || (len(parts) == 1 && !strings.HasSuffix(line, " "))

				var suggestions []string
				var currentWord string

				if isCommand {
					currentWord = line
					suggestions = cmdTrie.SearchPrefix(currentWord)
				} else {
					currentWord = parts[len(parts)-1]
					if strings.HasSuffix(line, " ") {
						currentWord = ""
					}
					suggestions = getFileSuggestions(currentWord)
				}
				if len(suggestions) == 1 {
					toAppend := suggestions[0][len(currentWord):]
					input = append(input, []rune(toAppend)...)
					fmt.Print(toAppend)
				} else if len(suggestions) > 1 {
					// Várias sugestões: imprime-as abaixo (como o bash)
					fmt.Print("\r\n")
					printCompact(suggestions)
					fmt.Printf("\r\ngosh | %s at %s > %s", usr.Username, pwd, string(input))
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

func getFileSuggestions(prefix string) []string {
	files, err := os.ReadDir(".")
	if err != nil {
		return nil
	}
	var result []string
	for _, file := range files {
		if strings.HasPrefix(file.Name(), prefix) {
			name := file.Name()
			if file.IsDir() {
				name += string("/")
			}
			result = append(result, name)
		}
	}
	return result
}

func printCompact(suggestions []string) {
	// Organiza em colunas simples para não poluir
	// Você pode evoluir isso para calcular a largura do terminal depois
	for i, s := range suggestions {
		fmt.Printf("%-20s", s) // Reserva 20 espaços por item
		if (i+1)%4 == 0 {      // Quebra linha a cada 4 itens
			fmt.Print("\r\n")
		}
	}
}
