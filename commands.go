package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dylandirosa98/blog-aggregator/internal/config"
	"github.com/dylandirosa98/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

type state struct {
	config *config.Config
	db     *database.Queries
}
type command struct {
	name string
	args []string
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("error not enough arguments: %v", cmd.args)
	}
	_, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		fmt.Printf("error getting user: %v", err)
		os.Exit(1)
	}
	err = s.config.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Printf("User has been set to: %v\n", cmd.args[0])
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		err := fmt.Errorf("error not enough arguments: %v", cmd.args)
		print(err.Error())
		return err
	}
	name := cmd.args[0]
	newContext := context.Background()
	newUser, err := s.db.CreateUser(
		newContext,
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      name,
		},
	)
	if err != nil {
		fmt.Printf("Error creating new user(possibly user already exists): %v\n", err)
		os.Exit(1)
	}
	s.config.SetUser(name)
	print("New user successfully created")
	log.Printf("user: %+v\n", newUser)
	return nil
}

func handlerResetUser(s *state, cmd command) error {
	if len(cmd.args) < 0 {
		err := fmt.Errorf("error not enough arguments: %v", cmd.args)
		print(err.Error())
		return err
	}
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		fmt.Printf("Error resetting users: %v\n", err)
		os.Exit(1)
	}
	return nil
}

func handlerGetUsers(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		err := fmt.Errorf("error too many arguments: %v", cmd.args)
		print(err.Error())
		return err
	}
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		fmt.Printf("Error getting users: %v\n", err)
		os.Exit(1)
	}
	for _, user := range users {
		if user.Name == s.config.Current_user_name {
			fmt.Printf("* %v (current)\n", user.Name)
		} else {
			fmt.Printf("* %v\n", user.Name)
		}
	}
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
