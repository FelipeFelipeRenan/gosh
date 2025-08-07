package builtin

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)
var ErrExit = fmt.Errorf("exit requested")


func Exec(args []string) (bool, error) {
	if len(args) == 0 {
		return false, nil
	}

	switch args[0] {
	case "cd":
		return true, cd(args)
	case "exit":
		return true, ErrExit
	default:
		return false, nil
	}
	//return false, nil
}

func cd(args []string) error {
	var dir string

	if len(args) < 2 {
		// sem argumentos, vai para home
		usr, err := user.Current()
		if err != nil {
			return fmt.Errorf("erro ao obter usuário atual: %v", err)
		}

		dir = usr.HomeDir

	} else {
		dir = args[1]
	}

	dir, err := filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("caminho inválido: %v", err)
	}

	err = os.Chdir(dir)
	if err != nil {
		return fmt.Errorf("não foi possivel mudar para %s: %v", dir, err)
	}

	return nil
}
