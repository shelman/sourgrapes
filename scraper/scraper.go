package scraper

import ()

var (
	fileDir = "/Users/sam/code/skunk/txt"
)

type Scraper interface {
	Scrape() error
}
