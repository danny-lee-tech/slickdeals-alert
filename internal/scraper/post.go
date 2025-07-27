package scraper

import (
	"fmt"
	"strings"
	"time"
)

type Post struct {
	Id         string
	Url        string
	Title      string
	Category   string
	ReplyCount int
	ViewCount  int
	Created    time.Time
	LastPosted time.Time
	Rank       int
	Reason     string
}

func (post *Post) ToString() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Id: %s\n", post.Id))
	sb.WriteString(fmt.Sprintf("Title: %s\n", post.Title))
	sb.WriteString(fmt.Sprintf("Rank: %d\n", post.Rank))
	sb.WriteString(fmt.Sprintf("Category: %s\n", post.Category))
	sb.WriteString(fmt.Sprintf("Replies: %d\n", post.ReplyCount))
	sb.WriteString(fmt.Sprintf("View Count: %d\n", post.ViewCount))
	sb.WriteString(fmt.Sprintf("Created: %s\n", post.Created.Format("2006-01-02 15:04:05")))
	sb.WriteString(fmt.Sprintf("Last Posted: %s\n", post.LastPosted.Format("2006-01-02 15:04:05")))
	sb.WriteString(fmt.Sprintf("URL: %s\n", post.Url))
	return sb.String()
}

func (post *Post) PrintableInfo() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Title: %s\n", post.Title))
	sb.WriteString(fmt.Sprintf("Reason: %s\n", post.Reason))
	sb.WriteString(fmt.Sprintf("Rank: %d\n", post.Rank))
	sb.WriteString(fmt.Sprintf("Category: %s\n", post.Category))
	sb.WriteString(fmt.Sprintf("Replies: %d\n", post.ReplyCount))
	sb.WriteString(fmt.Sprintf("View Count: %d\n", post.ViewCount))
	sb.WriteString(fmt.Sprintf("URL: %s\n", post.Url))
	return sb.String()
}
