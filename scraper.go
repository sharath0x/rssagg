package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sharath0x/rssagg/internal/database"
)

func StartScraping(
	db *database.Queries,
	concurrency int,
	timeBetweenRequest time.Duration,
) {
	log.Printf("Scraping on %v goroutines for %s duration", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Printf("error fetching feeds, %v", err)
			continue
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			ScrapeFeed(db, wg, feed)
		}
		wg.Wait()
	}
}

func ScrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	_, err := db.MarkFeedasFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("failed marking feeds as marked: ", err)
		return
	}

	rssFeed, err := UrlToFeed(feed.Url)
	if err != nil {
		log.Println("failed fetching feed: ", err)
		return
	}

	for _, item := range rssFeed.Channel.Item {

		desc := sql.NullString{}
		if item.Description != "" {
			desc.String = item.Description
			desc.Valid = true
		}

		pubAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("failed to parse date: %v , err: %v ", item.PubDate, err)
			continue
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Description: desc,
			PublishedAt: pubAt,
			Url:         item.Link,
			FeedID:      feed.ID,
		})

		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			} else {
				log.Println("failed to create post: ", err)
				continue
			}
		}
	}
	log.Printf("Feed %s collected, %v posts found\n", rssFeed.Channel.Title, len(rssFeed.Channel.Item))
}
