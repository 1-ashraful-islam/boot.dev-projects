package scrapper

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/1-ashraful-islam/boot.dev-projects/18.BlogAggregator/internal/database"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type FeedData struct {
	Items []struct {
		Title       string `xml:"title"`
		Link        string `xml:"link"`
		PubDate     string `xml:"pubDate"`
		Description string `xml:"description"`
	} `xml:"channel>item"`
}

func ScrapeFeed(ctx context.Context, db *database.Queries, feed database.Feed) error {
	// Scrape the feed
	// Save the feed to the database
	fmt.Println("Scraping feed", feed.Url)

	body, err := fetchURL(feed.Url)
	if err != nil {
		return errors.Wrap(err, "fetching feed failed for "+feed.Url)
	}

	feedData := &FeedData{}

	if err := xml.Unmarshal(body, &feedData); err != nil || len(feedData.Items) == 0 {
		return errors.Wrap(err, "parsing feed failed for "+feed.Url)
	}

	// update the last fetched at time
	_, err = db.MarkFeedAsFetched(ctx, database.MarkFeedAsFetchedParams{
		ID: feed.ID,
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return errors.Wrap(err, "updating feed to the database for "+feed.Url)
	}

	for _, item := range feedData.Items {

		// Check if the post already exists
		if _, err := db.GetPostByURL(ctx, item.Link); err == nil {
			continue
		}
		parsedTime, err := parseTime(item.PubDate)
		if err != nil {
			return errors.Wrap(err, "parsing published time failed for "+feed.Url)
		}

		_, err = db.CreatePost(ctx, database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			FeedID:      feed.ID,
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishDate: parsedTime,
		})
		if err != nil {
			return errors.Wrap(err, "creating post to the database for "+feed.Url)
		}
	}
	log.Println("Scraped feed", feed.Url, "with", len(feedData.Items), "items")

	return nil
}

func fetchURL(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func parseTime(s string) (time.Time, error) {
	var formats = []string{time.RFC1123Z, time.RFC1123}
	var parsedTime time.Time
	var err error

	for _, format := range formats {
		parsedTime, err = time.Parse(format, s)
		if err == nil {
			break
		}
	}

	if err != nil {
		return time.Time{}, errors.Wrap(err, " expected formats: "+strings.Join(formats, ", ")+" got: "+s)
	}
	return parsedTime, nil
}
