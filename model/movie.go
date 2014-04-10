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
	Id    bson.ObjectId `bson:"_id,omitempty"`
	Title string        `bson:"t"`
}

func (self *Movie) Insert() error {
	sess, d, err := db.NewSessionFactory().GetSession()
	if err != nil {
		return fmt.Errorf("couldn't get session: %v", err)
	}
	defer sess.Close()

	return d.C(MOVIE_COLLECTION).Insert(self)
}
