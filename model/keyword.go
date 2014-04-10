package model

import (
	"fmt"
	"labix.org/v2/mgo/bson"
	"shelman/sourgrapes/db"
)

const (
	KEYWORD_COLLECTION = "keywords"
)

type Keyword struct {
	Id     bson.ObjectId `bson:"_id,omitempty"`
	Word   string        `bson:"kw"`
	Movies []string      `bson:"mv"`
}

func FindKeywords(words []string) ([]Keyword, error) {
	sess, d, err := db.GetFactory().GetSession()
	if err != nil {
		return nil, fmt.Errorf("couldn't get session: %v", err)
	}
	defer sess.Close()

	keywords := []Keyword{}
	err = d.C(KEYWORD_COLLECTION).Find(
		bson.M{
			"kw": bson.M{
				"$in": words,
			},
		},
	).All(&keywords)
	return keywords, err
}

func (self *Keyword) Insert() error {
	sess, d, err := db.GetFactory().GetSession()
	if err != nil {
		return fmt.Errorf("couldn't get session: %v", err)
	}
	defer sess.Close()

	return d.C(KEYWORD_COLLECTION).Insert(self)
}
