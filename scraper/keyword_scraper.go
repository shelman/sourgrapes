package scraper

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"shelman/sourgrapes/model"
	"strings"
)

var (
	keywordFilename = "keywords.list"
	keywordHeader   = "8: THE KEYWORDS LIST"
	alphaNumeric    = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

type KeywordScraper struct{}

func (self *KeywordScraper) Scrape() error {

	// read in file
	fileBytes, err := ioutil.ReadFile(filepath.Join(fileDir, keywordFilename))
	if err != nil {
		return fmt.Errorf("error reading file %v: %v", keywordFilename, err)
	}

	// conv to string
	fileStr := string(fileBytes)

	// find where the keywords start, peel off the rest
	headerIdx := strings.Index(fileStr, keywordHeader)
	if headerIdx == -1 {
		return fmt.Errorf("keywords header not found")
	}
	postHeader := fileStr[headerIdx:]

	// split it into lines
	byLine := strings.Split(postHeader, "\n")

	// loop over lines (skip the header)
	var currMovie *model.Movie = nil
	keywords := map[string][]string{}
	fmt.Println(fmt.Sprintf("%v movies", len(byLine)))
	for idx, line := range byLine[1:] {
		if len(line) == 0 || !strings.Contains(alphaNumeric, line[:1]) ||
			strings.Contains(line, "(VG)") || strings.Contains(line, "(TV)") ||
			strings.Contains(line, "(V)") || strings.Contains(line, "/I)") ||
			strings.Contains(line, "{{SUSPENDED}}") {
			continue
		}

		// split the line into title + year + keyword
		lastParen := strings.LastIndex(line, ")")
		title := line[:lastParen+1]
		year, title := (title[lastParen-4 : lastParen]),
			strings.Trim(strings.ToLower(title[:lastParen-5]), " ")
		keyword := strings.Trim(line[lastParen+1:], "\t")

		if currMovie == nil {
			// first movie
			currMovie = &model.Movie{
				Title:    title,
				Year:     year,
				Keywords: []string{keyword},
			}
		} else if currMovie.Title == title {
			// still on the same movie
			currMovie.Keywords = append(currMovie.Keywords, keyword)
		} else {
			// new movie! insert the old and reset
			if err := currMovie.Insert(); err != nil {
				fmt.Println(fmt.Sprintf("error saving movie %v: %v",
					currMovie.Title, err))
			}
			currMovie = &model.Movie{
				Title:    title,
				Year:     year,
				Keywords: []string{keyword},
			}
		}

		keywords[keyword] = append(keywords[keyword], title)

		if idx%200 == 0 {
			fmt.Println(fmt.Sprintf("movies at %v", idx))
		}

	}

	// save the last one
	if currMovie != nil {
		if err := currMovie.Insert(); err != nil {
			fmt.Println(fmt.Sprintf("error saving movie %v: %v", currMovie.Title,
				err))
		}
	}

	fmt.Println("Done saving movies!")

	// save the keywords
	fmt.Println(fmt.Sprintf("%v keywords", len(keywords)))
	i := 0
	for k, mvs := range keywords {
		if err := (&model.Keyword{Word: k, Movies: mvs}).Insert(); err != nil {
			fmt.Println(fmt.Sprintf("error saving keyword %v: %v", k, err))
		}
		i++
		if i%1000 == 0 {
			fmt.Println(fmt.Sprintf("keywords at %v", i))
		}
	}

	return nil

}
