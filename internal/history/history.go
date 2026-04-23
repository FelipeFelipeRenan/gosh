package history

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var commands []string
var index int = -1

type History struct {
	entries  []string
	pos      int
	filePath string
}

func New(path string) *History {
	h := &History{
		filePath: path,
	}

	file, err := os.Open(path)
	if err == nil {
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				fmt.Printf("Erro ao fechar arquivo: %v\n", err)
			}
		}(file)
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" {
				h.entries = append(h.entries, line)
			}
		}
	}
	h.pos = len(h.entries)
	return h
}

func (h *History) Add(command string) {
	if command == "" {
		return
	}

	if len(h.entries) > 0 && h.entries[len(h.entries)-1] == command {
		h.pos = len(h.entries)
		return
	}
	h.entries = append(h.entries, command)
	h.pos = len(h.entries)

	if h.filePath != "" {
		file, err := os.OpenFile(h.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			defer func(file *os.File) {
				err := file.Close()
				if err != nil {
					fmt.Printf("Erro ao fechar arquivo: %v\n", err)
				}
			}(file)
			_, err2 := file.WriteString(command + "\n")
			if err2 != nil {
				return
			}
		}
	}
}

func (h *History) Prev() string {
	if len(h.entries) == 0 || h.pos <= 0 {
		return ""
	}
	h.pos--
	return h.entries[h.pos]
}

func (h *History) Next() string {
	if len(h.entries) == 0 || h.pos >= len(h.entries)-1 {
		h.pos = len(h.entries)
		return ""
	}

	h.pos++
	return h.entries[h.pos]
}

func (h *History) ResetPos() {
	h.pos = len(h.entries)
}
