package parser

import "strings"

func Parse(input string) []string{
	return strings.Fields(input)
}

