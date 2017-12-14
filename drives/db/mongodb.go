package db

import (
	"gopkg.in/mgo.v2"
	"log"
)

var MgoSession *mgo.Session = &mgo.Session{}

func init()  {
	var err error
	MgoSession, err = mgo.Dial("")
	if err != nil {
		log.Println("MongoDB：%v\n", err)
	} else if err = MgoSession.Ping(); err != nil {
		log.Println("MongoDB：%v\n", err)
	}
	MgoSession.SetMode(mgo.Monotonic, true)
	//default is 4096
	MgoSession.SetPoolLimit(300)
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

func CloneSession() *mgo.Session {
	return MgoSession.Clone()
}