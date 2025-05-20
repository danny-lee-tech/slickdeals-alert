package main

import (
	"fmt"
	"log"

	"github.com/danny-lee-tech/slickdeals-alert/internal/scraper"
)

func main() {
	scraper1 := scraper.Scraper{
		VoteFilter:        2,
		NotifyMinimumRank: 7,
	}

	results, err := scraper1.Execute()
	if err != nil {
		log.Fatal("Error:", err)
		return
	}

	if len(results) == 0 {
		fmt.Println("No results found")
		return
	}

	for i, result := range results {
		fmt.Printf("Result %d:\n", i)
		fmt.Printf("Title: %s\n", result.Text)
		fmt.Printf("URL: %s\n", result.Url)
		fmt.Printf("Rank: %d\n", result.Rank)
	}
}
