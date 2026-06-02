package main

import (
	"fmt"
	"os"

	"github.com/dylandirosa98/blog-aggregator/internal/config"
)

func main() {
	con, err := config.Read()
	if err != nil {
		print(fmt.Errorf("error reading config file: %v", err))
	}
	theState := state{config: &con}
	theCommands := &commands{mapCommands: make(map[string]func(*state, command) error)}
	theCommands.register("login", handlerLogin)
	if err != nil {
		print(fmt.Errorf("error registering command: %v", err))
	}
	if len(os.Args) < 2 {
		if len(os.Args) > 1 {
			fmt.Printf("not enough arguments\n")
		} else {
			fmt.Printf("a username is required\n")
		}
		os.Exit(1)
	}
	cmd := command{os.Args[1], os.Args[2:]}
	err = theCommands.run(&theState, cmd)
	if err != nil {
		print(fmt.Errorf("error running command: %v", err))
		os.Exit(1)
	}
	os.Exit(0)
}
