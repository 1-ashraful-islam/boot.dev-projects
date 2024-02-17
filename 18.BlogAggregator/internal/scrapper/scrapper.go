package scrapper

import (
	"context"
	"fmt"

	"github.com/1-ashraful-islam/boot.dev-projects/18.BlogAggregator/internal/database"
)

func ScrapeFeeds(ctx context.Context, db *database.Queries, feeds []database.Feed) {
	for _, feed := range feeds {
		// Scrape the feed
		// Save the feed to the database
		fmt.Println("Scraping feed", feed.Url)
	}
}
