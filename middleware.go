package main

import (
	"context"
	"fmt"

	"github.com/dylandirosa98/blog-aggregator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(s *state, cmd command) error {
	return func(s *state, cmd command) error {
		currentUser, err := s.db.GetUser(context.Background(), s.config.Current_user_name)
		if err != nil {
			fmt.Printf("get current user failed: %v\n", err)
			return err
		}
		return handler(s, cmd, currentUser)
	}
}
