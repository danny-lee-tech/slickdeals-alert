package scraper

import (
	"fmt"
	"strings"
)

type Post struct {
	Id    string
	Url   string
	Title string
	Rank  int
}

func (post *Post) ToString() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Id: %s\n", post.Id))
	sb.WriteString(fmt.Sprintf("Title: %s\n", post.Title))
	sb.WriteString(fmt.Sprintf("Rank: %d\n", post.Rank))
	sb.WriteString(fmt.Sprintf("URL: %s\n", post.Url))
	return sb.String()
}
