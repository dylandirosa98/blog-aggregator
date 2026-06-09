package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/dylandirosa98/blog-aggregator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.args) > 0 {
		n, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf("invalid limit: %w", err)
		}
		limit = n
	}
	theLimit := int32(limit)
	arg := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  theLimit,
	}
	posts, err := s.db.GetPostsForUser(context.Background(), arg)
	if err != nil {
		fmt.Printf("Error getting posts: %s\n", err)
		return err
	}
	for _, post := range posts {
		fmt.Printf("Title %s\n", post.Title)
		fmt.Printf("Url: %s\n", post.Url)
	}
	return nil
}
