package db

import (
    "gopkg.in/mgo.v2"
    "log"
)

//Mongo session Exported object
var MgoSession *mgo.Session

func init(){
    log.Println("Starting mongod session")
    session, err := mgo.Dial("127.0.0.1") 
    if err!=nil{
	panic(err)
    }
    MgoSession = session
}
