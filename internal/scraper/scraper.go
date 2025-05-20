package scraper

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

const (
	baseurl     string = "https://www.slickdeals.net/%s"
	url         string = "https://www.slickdeals.net/forums/filtered/?daysprune=7&vote=%d&f=9&sort=threadstarted&order=desc&r=1"
	ignoreTitle string = "A tl;dr of Slickdeals Rules and Guidelines and all that fun stuff"
)

type Scraper struct {
	VoteFilter        int // Search Filter on minimum number of votes. Used to determine the URL to scrape, specifically the vote query parameter
	NotifyMinimumRank int // the minimum number of thumbs up x 2 before a notification occurs
}

func (r Scraper) Execute() ([]Result, error) {
	htmlContent, err := r.scrape()
	if err != nil {
		return nil, err
	}

	selection, err := retrieveTableElement(htmlContent)
	if err != nil {
		return nil, err
	}

	return r.collect(selection)
}

func (r Scraper) getScrapeURL() string {
	return fmt.Sprintf(url, r.VoteFilter)
}

func (r Scraper) scrape() (string, error) {
	scrapeUrl := r.getScrapeURL()
	fmt.Printf("Scraping URL: %s\n", scrapeUrl)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	c, _ := chromedp.NewContext(ctx)
	defer chromedp.Cancel(c)

	var htmlContent string
	err := chromedp.Run(c,
		chromedp.Navigate(scrapeUrl),
		chromedp.WaitVisible("#threadbits_forum_9", chromedp.ByID),
		chromedp.OuterHTML("body", &htmlContent),
	)
	if err != nil {
		log.Fatal("Error:", err)
		return "", err
	}

	return htmlContent, nil
}

func retrieveTableElement(htmlContent string) (*goquery.Selection, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		log.Fatal("Goquery error:", err)
		return nil, err
	}

	selection := doc.Find("#threadbits_forum_9")
	return selection, nil
}

func (r Scraper) collect(selection *goquery.Selection) ([]Result, error) {
	var results []Result
	selection.Find("td[id^='td_threadtitle_'] .concat-thumbs").Each(func(index int, row *goquery.Selection) {
		class, _ := row.Attr("class")
		re, err := regexp.Compile(`rating(\d+)`)
		if err == nil {
			matches := re.FindStringSubmatch(class)
			rating, err := strconv.Atoi(matches[1])
			if err == nil && rating >= r.NotifyMinimumRank {
				threadElement := row.Parent().Parent().Parent()
				anchorElement := threadElement.Find("span.blueprint a").First()
				returnText := anchorElement.Text()

				if returnText == ignoreTitle {
					return
				}

				hrefValue, _ := anchorElement.Attr("href")
				returnUrl := fmt.Sprintf(baseurl, hrefValue)
				result := Result{
					Url:  returnUrl,
					Text: returnText,
					Rank: rating,
				}
				results = append(results, result)
			}
		}
	})

	return results, nil
}
