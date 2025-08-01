package main

import (
	"bufio"
	"fmt"

	"os"
	"strings"

	"github.com/FelipeFelipeRenan/gosh/internal/executor"
	"github.com/FelipeFelipeRenan/gosh/internal/parser"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for{
		fmt.Print("gosh> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("erro ao ler entrada: ", err)
			break
		}

		input = strings.TrimSpace(input)
		if input == ""{
			continue
		}

		args := parser.Parse(input)

		if err := executor.Exec(args); err != nil{
			fmt.Println("erro ao executar comando: ", err)
		}

	}
	
}