package main

import (
	"fmt"
	"log"

	"github.com/danny-lee-tech/slickdeals-alert/internal/scraper"
)

func main() {
	scraper1 := scraper.Scraper{
		VoteFilter:        2,
		NotifyMinimumRank: 5,
	}

	results, err := scraper1.Execute()
	if err != nil {
		log.Fatal("Error:", err)
		return
	}

	for i, result := range results {
		fmt.Printf("Result %d:\n", i)
		fmt.Println(result.Text)
		fmt.Println(result.Url)
	}
}
