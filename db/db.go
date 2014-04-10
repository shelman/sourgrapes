package db

import (
	"labix.org/v2/mgo"
	"time"
)

var (
	url         = "mongodb://localhost:27017"
	dbName      = "sourgrapes"
	dialTimeout = time.Second * 5
)

type SessionFactory struct {
	Url           string
	DBName        string
	DialTimeout   time.Duration
	masterSession *mgo.Session
}

func NewSessionFactory() *SessionFactory {
	return &SessionFactory{
		Url:         url,
		DBName:      dbName,
		DialTimeout: dialTimeout,
	}
}

func (self *SessionFactory) GetSession() (*mgo.Session, *mgo.Database, error) {

	// init the master session if necessary
	if self.masterSession == nil {
		var err error
		self.masterSession, err = mgo.DialWithTimeout(self.Url,
			self.DialTimeout)
		if err != nil {
			return nil, nil, err
		}
	}

	// copy the master session
	sessionCopy := self.masterSession.Copy()
	return sessionCopy, sessionCopy.DB(self.DBName), nil

}
