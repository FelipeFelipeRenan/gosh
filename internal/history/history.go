package history


var commands []string
var index int = -1

func Add(command string)  {
	if command != ""{
		commands = append(commands, command)
		index = len(command) // sempre depois do ultimo
	}
}

func Prev() string  {
	if len(commands) == 0 || index <= 0{
		return ""
	}
	index--
	return commands[index]
}

func Next() string{
	if len(commands) == 0 || index >= len(commands)-1{
		index = len(commands)
		return ""
	}

	index++
	return commands[index]
}
