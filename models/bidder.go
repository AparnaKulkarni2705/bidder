package models

import (
    mgo "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"

    "log"
)

const BIDDER_COLLECTION_NAME = "bidder"

type Bidder struct {
    ID	    bson.ObjectId   `bson : "_id" json : "id"`
    Name    string	    `bson: "name" json : "name"`
    Email   string	    `bson: "email" json : "email"`
}

func (bidder *Bidder) Insert(session *mgo.Session) error {
    bidder_collection := session.DB("di_bidder_db").C(BIDDER_COLLECTION_NAME)
    bidder.ID = bson.NewObjectId()

    err:=bidder_collection.Insert(bidder)
    if err !=nil{
	log.Fatalf("Error inserting bidder %s: Error: %s", bidder.Email, err)
    }
    return nil
}
