package main

import (
	"fmt"

	"github.com/mmcdole/gofeed"
)

func main() {
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL("https://g1.globo.com/rss/g1/")
	for _, item := range feed.Items {
		fmt.Println("Título:", item.Title)
		fmt.Println("Link:", item.Link)
	}
}
