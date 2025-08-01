package main

import (
	"bufio"
	"fmt"
	"os/user"

	"os"
	"strings"

	"github.com/FelipeFelipeRenan/gosh/internal/builtin"
	"github.com/FelipeFelipeRenan/gosh/internal/executor"
	"github.com/FelipeFelipeRenan/gosh/internal/parser"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	usr, err := user.Current()
	if err != nil {
		fmt.Println("erro de usuário: ", err)
	}

	for {
		fmt.Printf("gosh | %s > ", usr.Username)

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("erro ao ler entrada: ", err)
			break
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		args := parser.Parse(input)

		handled, err := builtin.Exec(args)
		if err != nil {
			fmt.Println("erro no comando interno: ", err)
			continue
		}
		if handled {
			continue // comando interno executado, pula executor externo
		}

		if err := executor.Exec(args); err != nil {
			fmt.Println("erro ao executar comando: ", err)
		}

	}

}
