package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/danny-lee-tech/slickdeals-alert/internal/emailer"
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

	var sb strings.Builder

	for _, result := range results {
		sb.WriteString(fmt.Sprintf("Title: %s\n", result.Text))
		sb.WriteString(fmt.Sprintf("URL: %s\n", result.Url))
		sb.WriteString(fmt.Sprintf("Rank: %d\n\n", result.Rank))
	}

	fmt.Print(sb.String())
	emailer.Email(sb.String())
}
