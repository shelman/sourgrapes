package main

import (
	"fmt"
	"shelman/sourgrapes/model"
)

func main() {
	mv, err := model.FindMovie("jaws")
	if err != nil {
		panic(fmt.Sprintf("error finding movie: %v", err))
	}

	err = mv.GetMoviesWithSameKeywords()
	if err != nil {
		panic(fmt.Sprintf("error finding movies with same keywords: %v", err))
	}
}
