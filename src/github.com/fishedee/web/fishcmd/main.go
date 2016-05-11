package main

import (
	"fmt"
	"github.com/fishedee/web/fishcmd/command"
	"os"
)

type commandHandlerType func(argv []string) (string, error)

func main() {
	args := os.Args
	args = args[1:]
	cmd := ""
	if len(args) != 0 {
		cmd = args[0]
		args = args[1:]
	}

	commandHandler := []struct {
		name    string
		handler commandHandlerType
	}{
		{"help", command.Help},
		{"version", command.Version},
		{"run", command.Run},
		{"test", command.Test},
	}

	var singleCommandHandler commandHandlerType
	for _, singleCommand := range commandHandler {
		if singleCommand.name == cmd {
			singleCommandHandler = singleCommand.handler
			break
		}
	}
	if singleCommandHandler == nil {
		singleCommandHandler = commandHandler[0].handler
	}
	result, err := singleCommandHandler(args)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	} else {
		fmt.Println(result)
	}
}
