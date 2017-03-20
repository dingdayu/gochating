package db

import (
	"gopkg.in/mgo.v2"
	"log"
)

var MgoSession *mgo.Session = &mgo.Session{}

func init()  {
	MgoSession, err := mgo.Dial("")
	if err != nil {
		log.Println("MongoDB：%v\n", err)
	} else if err = MgoSession.Ping(); err != nil {
		log.Println("MongoDB：%v\n", err)
	}
}

func GetSession() *mgo.Session  {
	MgoSession, err := mgo.Dial("")
	if err != nil {
		log.Println("MongoDB：%v\n", err)
	} else if err = MgoSession.Ping(); err != nil {
		log.Println("MongoDB：%v\n", err)
	}
	return MgoSession
}