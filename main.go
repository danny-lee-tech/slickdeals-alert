package main

import (
	"fmt"
	"log"
	"os"

	"github.com/danny-lee-tech/slickdeals-alert/internal/emailer"
	"github.com/danny-lee-tech/slickdeals-alert/internal/scraper"
)

func main() {
	emailEnabled := true
	var emailPassword string
	if len(os.Args) <= 1 {
		emailEnabled = false
	} else {
		emailPassword = os.Args[1]
	}

	fmt.Println("Starting scraper")

	scraper1 := scraper.Scraper{
		VoteFilter:        1,
		NotifyMinimumRank: 8,
		GmailSetting: emailer.GmailSettingConfig{
			Enabled:            emailEnabled,
			SourceEmailAddress: "purewhiteasian@gmail.com",
			TargetEmailAddress: "onfire_22043@yahoo.com",
			Password:           emailPassword,
		},
	}

	err := scraper1.Execute()
	if err != nil {
		log.Fatal("Error:", err)
		return
	}
}
