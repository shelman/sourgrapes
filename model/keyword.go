package model

import (
	"fmt"
	"labix.org/v2/mgo/bson"
	"math/rand"
	"shelman/sourgrapes/db"
)

const (
	KEYWORD_COLLECTION = "keywords"
)

type Keyword struct {
	Id     bson.ObjectId `bson:"_id,omitempty" json:"_id"`
	Word   string        `bson:"kw" json:"word"`
	Movies []string      `bson:"mv" json:"movies"`
}

func GetRandomKeywords(howMany int) ([]Keyword, error) {
	sess, d, err := db.GetFactory().GetSession()
	if err != nil {
		return nil, fmt.Errorf("couldn't get session: %v", err)
	}
	defer sess.Close()

	total, err := d.C(KEYWORD_COLLECTION).Count()
	if err != nil {
		return nil, fmt.Errorf("error counting keywords: %v", err)
	}

	toRet := []Keyword{}
	for i := 0; i < howMany; i++ {
		toSkip := rand.Int() % total
		out := &Keyword{}
		err = d.C(KEYWORD_COLLECTION).Find(bson.M{}).Select(bson.M{"kw": 1}).
			Skip(toSkip).One(out)
		if err != nil {
			return nil, fmt.Errorf("error finding keyword: %v", err)
		}
		toRet = append(toRet, *out)
	}

	return toRet, nil
}

func FindKeyword(word string) (*Keyword, error) {
	sess, d, err := db.GetFactory().GetSession()
	if err != nil {
		return nil, fmt.Errorf("couldn't get session: %v", err)
	}
	defer sess.Close()

	keyword := &Keyword{}
	err = d.C(KEYWORD_COLLECTION).Find(bson.M{"kw": word}).One(keyword)
	return keyword, err
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
