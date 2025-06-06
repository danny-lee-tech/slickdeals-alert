package scraper

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

const baseurl string = "https://www.slickdeals.net/%s"

var ratingRegex = regexp.MustCompile(`rating(\d+)`)

func ConvertFromSelection(selection *goquery.Selection) Post {
	post := Post{}

	id, title, url := getIdTitleAndUrl(selection)
	post.Id = id
	post.Title = title
	post.Url = url
	post.Rank = getRank(selection)

	return post
}

func getIdTitleAndUrl(selection *goquery.Selection) (string, string, string) {
	anchorElement := selection.Find("span.blueprint a").First()
	title := anchorElement.Text()
	relativeUrl, _ := anchorElement.Attr("href")
	url := fmt.Sprintf(baseurl, relativeUrl)
	id, _ := anchorElement.Attr("id")

	return id, title, url
}

func getRank(selection *goquery.Selection) int {
	thumbElement := selection.Find(".concat-thumbs").First()
	class, _ := thumbElement.Attr("class")
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
