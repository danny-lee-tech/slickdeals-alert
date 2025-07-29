package scraper

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const baseurl string = "https://www.slickdeals.net/%s"

var (
	timeLocOnce sync.Once
	timeLoc     *time.Location
	ratingRegex = regexp.MustCompile(`rating(\d+)`)
)

func ConvertFromSelection(selection *goquery.Selection) Post {
	post := Post{}

	selection.Find("td").Each(func(index int, cell *goquery.Selection) {
		switch index {
		// Category = 1
		case 1:
			post.Category = getCategory(cell)
		// Id, Title, URL = 2
		case 2:
			id, title, url := getIdTitleAndUrl(cell)
			post.Id = id
			post.Title = title
			post.Url = url
			post.Rank = getRank(cell)
		// Thread Started = 3
		case 3:
			post.Created = getDate(cell)
		// Reply Count = 4
		case 4:
			post.ReplyCount = getReplyCount(cell)
		// View Count = 5
		case 5:
			post.ViewCount = getViewCount(cell)
		// Last Post = 6
		case 6:
			post.LastPosted = getDate(cell)
		}
	})

	return post
}

func getDate(selection *goquery.Selection) time.Time {
	element := selection.Find(".smallfont").First()
	dateText := strings.TrimSpace(element.Text())
	var day, year int
	var month time.Month
	if strings.Contains(dateText, "Today") {
		today := time.Now()
		month = today.Month()
		day = today.Day()
		year = today.Year()
	}

	timeElement := element.Find(".time").First()
	timeString := strings.TrimSpace(timeElement.Text())
	timeObject, _ := time.Parse("15:04 AM", timeString)

	loc, err := getLosAngelesTimeLocation()
	if err != nil {
		var t time.Time
		return t
	}
	return time.Date(year, month, day, timeObject.Hour(), timeObject.Minute(), 0, 0, loc)
}

func getLosAngelesTimeLocation() (*time.Location, error) {
	var loadErr error
	timeLocOnce.Do(func() {
		timeLoc, loadErr = time.LoadLocation("America/Los_Angeles")
	})
	if loadErr != nil {
		return nil, loadErr
	}
	return timeLoc, nil
}

func getCategory(selection *goquery.Selection) string {
	element := selection.Find(".threadCategoryForm button").First()
	category := element.Text()
	return category
}

func getReplyCount(selection *goquery.Selection) int {
	element := selection.Find("a").First()
	replyCountString := strings.ReplaceAll(element.Text(), ",", "")
	replyCount, _ := strconv.Atoi(replyCountString)
	return replyCount
}

func getViewCount(selection *goquery.Selection) int {
	viewCount, _ := strconv.Atoi(strings.ReplaceAll(strings.TrimSpace(selection.Text()), ",", ""))
	return viewCount
}

func getIdTitleAndUrl(selection *goquery.Selection) (string, string, string) {
	element := selection.Find("span.blueprint a").First()
	title := element.Text()
	relativeUrl, _ := element.Attr("href")
	url := fmt.Sprintf(baseurl, relativeUrl)
	id, _ := element.Attr("id")

	return id, title, url
}

func getRank(selection *goquery.Selection) int {
	element := selection.Find(".concat-thumbs").First()
	class, _ := element.Attr("class")
	matches := ratingRegex.FindStringSubmatch(class)
	rank := -1
	if len(matches) >= 1 {
		convertedRank, err := strconv.Atoi(matches[1])
		if err == nil {
			rank = convertedRank
		}
	}

	return rank
}
