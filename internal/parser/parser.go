package parser

import "strings"

func Parse(input string) []string {

	var args []string
	var currentArg strings.Builder

	inSingleQuote := false
	inDoubleQuote := false

	runes := []rune(strings.TrimSpace(input))

	for _, r := range runes {
		if inSingleQuote {
			if r == '\'' {
				inSingleQuote = false
			} else {
				currentArg.WriteRune(r)
			}
			continue
		}
		if inDoubleQuote {
			if r == '"' {
				inDoubleQuote = false
			} else {
				currentArg.WriteRune(r)
			}
			continue
		}

		switch r {
		case '\'':
			inSingleQuote = true
		case '"':
			inDoubleQuote = true
		case ' ', '\t':
			if currentArg.Len() > 0 {
				args = append(args, currentArg.String())
				currentArg.Reset()
			}
		default:
			currentArg.WriteRune(r)
		}
	}
	if currentArg.Len() > 0 {
		args = append(args, currentArg.String())
	}
	return args
}
