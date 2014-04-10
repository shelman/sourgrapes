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
	Id       bson.ObjectId `bson:"_id,omitempty"`
	Title    string        `bson:"t"`
	Year     string        `bson:"y"`
	Keywords []string      `bson:"kw"`
}

func (self *Movie) Insert() error {
	sess, d, err := db.GetFactory().GetSession()
	if err != nil {
		return fmt.Errorf("couldn't get session: %v", err)
	}
	defer sess.Close()

	return d.C(MOVIE_COLLECTION).Insert(self)
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

	fmt.Println(fmt.Sprintf("Main movie: %v", self.Title))
	for _, keyword := range keywords {
		fmt.Println(fmt.Sprintf("   Keyword: %v", keyword.Word))
		for _, mv := range keyword.Movies {
			fmt.Println(fmt.Sprintf("       %v", mv))
		}
	}

	return nil

}
