package scraper

import (
	"context"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/danny-lee-tech/slickdeals-alert/internal/emailer"
)

const (
	url         string = "https://www.slickdeals.net/forums/filtered/?daysprune=7&vote=%d&f=9&sort=threadstarted&order=desc&r=1"
	ignoreTitle string = "A tl;dr of Slickdeals Rules and Guidelines and all that fun stuff"
)

type Scraper struct {
	VoteFilter        int              // Search Filter on minimum number of votes. Used to determine the URL to scrape, specifically the vote query parameter
	NotifyMinimumRank int              // the minimum number of thumbs up x 2 before a notification occurs
	Emailer           *emailer.Emailer // Email Settings
}

func (r Scraper) Execute() error {
	htmlContent, err := r.scrape()
	if err != nil {
		return err
	}

	selection, err := retrieveTableElement(htmlContent)
	if err != nil {
		return err
	}

	posts, err := r.collect(selection)
	if err != nil {
		return err
	}

	if len(posts) == 0 {
		fmt.Println("No results found")
		return nil
	}

	postIds := getIds(posts)
	postsString := formatPosts(posts)

	matches, err := matchesPreviousPosts(postIds)
	if err != nil {
		return err
	}
	if matches {
		fmt.Println("Duplicate results. Avoiding email notification.")
		return nil
	}

	if r.Emailer != nil {
		r.Emailer.Email(postsString)
	} else {
		fmt.Println("Email notifications have been disabled")
	}

	return nil
}

func getIds(posts []Post) []string {
	var ids []string
	for _, post := range posts {
		ids = append(ids, post.Id)
	}
	return ids
}

func matchesPreviousPosts(postIds []string) (bool, error) {
	file, err := os.OpenFile("last_posts.txt", os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return false, err
	}
	defer file.Close()

	lastPostIdsBytes, err := io.ReadAll(file)
	if err != nil {
		return false, err
	}
	lastPostsString := string(lastPostIdsBytes)
	lastPostIds := strings.Split(lastPostsString, ",")

	matches := true
	for _, postId := range postIds {
		if !slices.Contains(lastPostIds, postId) {
			matches = false
		}
	}

	if matches {
		return true, nil
	}

	file, err = os.OpenFile("last_posts.txt", os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return false, err
	}
	defer file.Close()

	_, err = file.Write([]byte(strings.Join(postIds, ",")))
	if err != nil {
		return false, err
	}

	return false, nil
}

func formatPosts(posts []Post) string {
	var sb strings.Builder

	for _, post := range posts {
		sb.WriteString(post.PrintableInfo())
		sb.WriteString("\n")
	}

	fmt.Print(sb.String())
	return sb.String()
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
		return "", err
	}

	return htmlContent, nil
}

func retrieveTableElement(htmlContent string) (*goquery.Selection, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}

	selection := doc.Find("#threadbits_forum_9")
	return selection, nil
}

func (r Scraper) collect(selection *goquery.Selection) ([]Post, error) {
	var posts []Post
	selection.Find("tr[id^='sdpostrow']").Each(func(index int, row *goquery.Selection) {
		post := ConvertFromSelection(row)
		if post.Rank >= r.NotifyMinimumRank && post.Title != ignoreTitle {
			posts = append(posts, post)
		}
	})

	return posts, nil
}
