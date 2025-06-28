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
	"github.com/danny-lee-tech/slickdeals-alert/internal/pushbulleter"
)

const (
	url         string = "https://www.slickdeals.net/forums/filtered/?daysprune=7&vote=%d&f=9&sort=threadstarted&order=desc&r=1"
	ignoreTitle string = "A tl;dr of Slickdeals Rules and Guidelines and all that fun stuff"
)

type Scraper struct {
	VoteFilter        int // Search Filter on minimum number of votes. Used to determine the URL to scrape, specifically the vote query parameter
	NotifyMinimumRank int // the minimum number of thumbs up x 2 before a notification occurs
	Emailer           *emailer.Emailer
	PushBulleter      *pushbulleter.PushBulleter
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

	dedupPosts, err := deDuplicatePosts(posts)
	if err != nil {
		return err
	}
	postsString := formatPosts(dedupPosts)

	if len(dedupPosts) > 0 {
		if r.Emailer != nil {
			err = r.Emailer.Email(postsString)
			if err != nil {
				fmt.Println("Warning: Email", err)
			}
		} else {
			fmt.Println("Email notifications have been disabled")
		}

		if r.PushBulleter != nil {
			err = r.PushBulleter.PostToChannel(postsString)
			if err != nil {
				fmt.Println("Warning: PushBullet", err)
			}
		} else {
			fmt.Println("PushBullet pushes have been disabled")
		}
	} else {
		fmt.Println("Duplicate results. Avoiding email notification.")
	}

	return nil
}

func deDuplicatePosts(posts []Post) ([]Post, error) {
	file, err := os.OpenFile("last_posts.txt", os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lastPostIdsBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	lastPostsString := string(lastPostIdsBytes)
	lastPostIds := strings.Split(lastPostsString, ",")

	dedupedPosts := []Post{}
	for _, post := range posts {
		if !slices.Contains(lastPostIds, post.Id) {
			dedupedPosts = append(dedupedPosts, post)
			lastPostIds = append(lastPostIds, post.Id)
		}
	}

	if len(dedupedPosts) == 0 {
		return nil, nil
	}

	file, err = os.OpenFile("last_posts.txt", os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if len(lastPostIds) > 20 {
		elementsToDelete := len(lastPostIds) - 10
		lastPostIds = lastPostIds[elementsToDelete:]
	}
	_, err = file.Write([]byte(strings.Join(lastPostIds, ",")))
	if err != nil {
		return nil, err
	}

	return dedupedPosts, nil
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
