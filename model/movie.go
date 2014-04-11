package model

import (
	"fmt"
	"labix.org/v2/mgo/bson"
	"shelman/sourgrapes/db"
)

const (
	MOVIE_COLLECTION = "movies"
)

type Movie struct {
	Id       bson.ObjectId `bson:"_id,omitempty" json:"_id"`
	Title    string        `bson:"t" json:"title"`
	Year     string        `bson:"y" json:"year"`
	Keywords []string      `bson:"kw" json:"keywords"`
}

func (self *Movie) Insert() error {
	sess, d, err := db.GetFactory().GetSession()
	if err != nil {
		return fmt.Errorf("couldn't get session: %v", err)
	}
	defer sess.Close()

	return d.C(MOVIE_COLLECTION).Insert(self)
}

func FindMoviesWithPrefix(prefix string) ([]Movie, error) {
	sess, d, err := db.GetFactory().GetSession()
	if err != nil {
		return nil, fmt.Errorf("couldn't get session: %v", err)
	}
	defer sess.Close()

	movies := []Movie{}
	err = d.C(MOVIE_COLLECTION).Find(
		bson.M{
			"t": bson.RegEx{
				Pattern: "^" + prefix,
				Options: "",
			},
		},
	).Limit(20).All(&movies)
	return movies, err
}

func FindMovie(title string) (*Movie, error) {
	sess, d, err := db.GetFactory().GetSession()
	if err != nil {
		return nil, fmt.Errorf("couldn't get session: %v", err)
	}
	defer sess.Close()

	movie := &Movie{}
	err = d.C(MOVIE_COLLECTION).Find(
		bson.M{
			"t": title,
		},
	).One(movie)
	return movie, err
}

func FindMovies(titles []string) ([]Movie, error) {
	sess, d, err := db.GetFactory().GetSession()
	if err != nil {
		return nil, fmt.Errorf("couldn't get session: %v", err)
	}
	defer sess.Close()

	movies := []Movie{}
	err = d.C(MOVIE_COLLECTION).Find(
		bson.M{
			"t": bson.M{
				"$in": titles,
			},
		},
	).All(&movies)
	return movies, err

}

func (self *Movie) GetMoviesWithSameKeywords() error {

	// get the keywords from the db
	keywords, err := FindKeywords(self.Keywords)
	if err != nil {
		return fmt.Errorf("error getting keywords for movie %v: %v", self.Title,
			err)
	}

	simMap := map[string]int{}

	fmt.Println(fmt.Sprintf("Main movie: %v", self.Title))
	for _, keyword := range keywords {
		fmt.Println(fmt.Sprintf("   Keyword: %v", keyword.Word))
		for _, mv := range keyword.Movies {
			fmt.Println(fmt.Sprintf("       %v", mv))
			simMap[mv]++
		}
	}

	fmt.Println("Similarity scores")
	for mv, sim := range simMap {
		fmt.Println(fmt.Sprintf("%v: %v", mv, sim))
	}

	return nil

}
