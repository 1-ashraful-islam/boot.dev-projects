package scrapper

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/1-ashraful-islam/boot.dev-projects/18.BlogAggregator/internal/database"
)

type FeedData struct {
	Items []struct {
		Title string `xml:"title"`
		Link  string `xml:"link"`
	} `xml:"channel>item"`
}

func ScrapeFeeds(ctx context.Context, db *database.Queries, feeds []database.Feed) {
	for _, feed := range feeds {
		// Scrape the feed
		// Save the feed to the database
		fmt.Println("Scraping feed", feed.Url)

		body, err := fetchURL(feed.Url)
		if err != nil {
			fmt.Println("Error fetching feed", feed.Url, err)
			continue
		}

		feedData := &FeedData{}

		if err := xml.Unmarshal(body, &feedData); err != nil || len(feedData.Items) == 0 {
			fmt.Println("Error parsing feed", feed.Url, err)
			continue
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
			log.Println("Error updating feed to the database", feed.ID, err)
		}

		for _, item := range feedData.Items {
			fmt.Printf("Found title: %s\n  url: %s\n\n", item.Title, item.Link)
		}

	}
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
