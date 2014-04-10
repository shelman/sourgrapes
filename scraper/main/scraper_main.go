package main

import (
	"fmt"
	"shelman/sourgrapes/scraper"
)

func main() {

	// scrape keywords
	kwScraper := &scraper.KeywordScraper{}
	err := kwScraper.Scrape()
	if err != nil {
		fmt.Println(fmt.Sprintf("error scraping keywords: %v", err))
	}

}
