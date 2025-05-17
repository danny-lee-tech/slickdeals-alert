package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

const (
	baseurl string = "https://www.slickdeals.net/"
	url     string = "https://www.slickdeals.net/forums/filtered/?daysprune=7&vote=1&f=9&sort=threadstarted&order=desc&r=1"
)

func scrapeUrl() (string, error) {
	fmt.Println("Scraping URL: " + url)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	c, _ := chromedp.NewContext(ctx)
	defer chromedp.Cancel(c)

	var htmlContent string
	err := chromedp.Run(c,
		chromedp.Navigate(url),
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

func logTableElement(selection *goquery.Selection) {
	body, _ := selection.Html()
	file, err := os.Create("body.xml")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(body)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func main() {
	htmlContent, err := scrapeUrl()
	if err != nil {
		return
	}

	selection, err := retrieveTableElement(htmlContent)
	if err != nil {
		return
	}

	logTableElement(selection)

	file, err := os.Create("row.xml")
	if err != nil {
		log.Fatal(err)
		return
	}

	_, err = file.WriteString("<html>")
	if err != nil {
		log.Fatal(err)
		return
	}

	selection.Find("td[id^='td_threadtitle_'] .concat-thumbs").Each(func(index int, row *goquery.Selection) {
		class, _ := row.Attr("class")
		re, err := regexp.Compile(`rating(\d+)`)
		if err == nil {
			matches := re.FindStringSubmatch(class)
			rating, err := strconv.ParseInt(matches[1], 10, 64)
			if err == nil && rating > 5 {
				threadElement := row.Parent().Parent().Parent()
				anchorElement := threadElement.Find("span.blueprint a").First()
				text := anchorElement.Text()
				fmt.Println(text)
				hrefValue, _ := anchorElement.Attr("href")
				fmt.Println(baseurl + hrefValue)
			}
		}
		html, _ := row.Parent().Parent().Parent().Html()
		_, err = file.WriteString(html)
		if err != nil {
			log.Fatal(err)
			return
		}
	})

	_, err = file.WriteString("</html>")
	if err != nil {
		log.Fatal(err)
		return
	}
}
