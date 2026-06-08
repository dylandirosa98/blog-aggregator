package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/dylandirosa98/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		fmt.Printf("fetch feed request failed: %v\n", err)
		return nil, err
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("fetch feed response failed: %v\n", err)
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("read all failed: %v\n", err)
		return nil, err
	}
	feed := &RSSFeed{}
	if err := xml.Unmarshal(body, feed); err != nil {
		fmt.Printf("unmarshal failed: %v\n", err)
		return nil, err
	}
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}
	return feed, nil
}

func handlerAgg(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		fmt.Printf("fetch feed failed: %v\n", err)
		return err
	}
	fmt.Printf("%+v\n", *feed)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		fmt.Printf("invalid arguments\n")
		os.Exit(1)
	}
	currentUser, err := s.db.GetUser(context.Background(), s.config.Current_user_name)
	if err != nil {
		fmt.Printf("get current user failed: %v\n", err)
		return err
	}
	userID := currentUser.ID
	arg := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    userID,
	}
	_, err = s.db.CreateFeed(context.Background(), arg)
	if err != nil {
		fmt.Printf("create feed failed: %v\n", err)
		return err
	}
	feedFollowArg := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    arg.ID,
	}
	_, err = s.db.CreateFeedFollow(context.Background(), feedFollowArg)
	if err != nil {
		fmt.Printf("create feed follow failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%+v\n", arg)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		fmt.Printf("invalid arguments\n")
		os.Exit(1)
	}
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		fmt.Printf("get feeds failed: %v\n", err)
		os.Exit(1)
	}
	for i, feed := range feeds {
		fmt.Printf("Feed #%d: %v\n", i, feed.FeedName)
		fmt.Printf("Url: %v\n", feed.Url)
		fmt.Printf("User Name: %v\n", feed.UserName)
	}
	return nil
}

func handlerFollow(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		fmt.Printf("invalid arguments\n")
		os.Exit(1)
	}
	user, err := s.db.GetUser(context.Background(), s.config.Current_user_name)
	if err != nil {
		fmt.Printf("get current user failed: %v\n", err)
		os.Exit(1)
	}
	feed, err := s.db.QueryFeed(context.Background(), cmd.args[0])
	if err != nil {
		fmt.Printf("get feed failed: %v\n", err)
		os.Exit(1)
	}
	arg := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	_, err = s.db.CreateFeedFollow(context.Background(), arg)
	if err != nil {
		fmt.Printf("create feed follow failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Feed: %v\nUser: %v\n", feed.Name, cmd.name)
	return nil
}

func handlerFollowing(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		fmt.Printf("invalid arguments\n")
		os.Exit(1)
	}
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), s.config.Current_user_name)
	if err != nil {
		fmt.Printf("get feed follows failed: %v\n", err)
		os.Exit(1)
	}
	for i, _ := range feeds {
		fmt.Printf("%v\n", feeds[i].Name)
	}
	return nil
}
