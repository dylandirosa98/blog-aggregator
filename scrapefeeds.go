package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/dylandirosa98/blog-aggregator/internal/database"
	"github.com/google/uuid"
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
	Feed := *feed
	for i := range Feed.Channel.Item {
		formats := []string{time.RFC1123Z, time.RFC1123, time.RFC3339}
		var t time.Time
		var err error
		for _, format := range formats {
			t, err = time.Parse(format, Feed.Channel.Item[i].PubDate)
			if err == nil {
				break
			}
		}
		description := sql.NullString{
			String: Feed.Channel.Item[i].Description,
			Valid:  true,
		}
		if Feed.Channel.Item[i].Description == "" {
			description.Valid = false
		}
		publishedAt := sql.NullTime{
			Time:  t,
			Valid: true,
		}
		if t == (time.Time{}) {
			publishedAt.Valid = false
		}
		id := uuid.NullUUID{
			UUID:  nextFeed.ID,
			Valid: true,
		}
		arg := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       Feed.Channel.Item[i].Title,
			Url:         Feed.Channel.Item[i].Link,
			Description: description,
			PublishedAt: publishedAt,
			FeedID:      id,
		}
		err = s.db.CreatePost(context.Background(), arg)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				continue
			}
			log.Printf("Error creating post: %v", err)
		}
	}
	return nil
}
