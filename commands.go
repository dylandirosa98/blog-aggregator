package main

import (
	"fmt"

	"github.com/dylandirosa98/blog-aggregator/internal/config"
)

type state struct {
	config *config.Config
}
type command struct {
	name string
	args []string
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("error not enough arguments: %v", cmd.args)
	}
	err := s.config.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Printf("User has been set to: %v\n", cmd.args[0])
	return nil
}

type commands struct {
	mapCommands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	theFunc, ok := c.mapCommands[cmd.name]
	if !ok {
		fmt.Printf("Unknown command: %v\n", cmd.name)
		return fmt.Errorf("unknown command: %v", cmd.name)
	}
	err := theFunc(s, cmd)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.mapCommands[name] = f
}
