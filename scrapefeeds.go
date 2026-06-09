package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/dylandirosa98/blog-aggregator/internal/database"
)

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		fmt.Print("Error fetching next feed")
		os.Exit(1)
	}
	myNullableTime := sql.NullTime{
		Time:  time.Now(),
		Valid: true, // This tells the database "Yes, this is a real time, not NULL"
	}
	arg := database.MarkFeedFetchedParams{
		LastFetchedAt: myNullableTime,
		ID:            nextFeed.ID,
	}
	err = s.db.MarkFeedFetched(context.Background(), arg)
	if err != nil {
		fmt.Print("Error fetching next feed")
		os.Exit(1)
	}
	feed, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		fmt.Print("Error fetching next feed")
		os.Exit(1)
	}
	fmt.Printf("%v\n", feed.Channel.Title)
	Feed := *feed
	for i := range Feed.Channel.Item {
		fmt.Printf("%v\n", Feed.Channel.Item[i].Title)
	}
	return nil
}
