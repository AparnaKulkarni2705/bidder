package models

import "gopkg.in/mgo.v2/bson"

type Bidder struct {
    ID	    bson.ObjectID   `bson : "_id" json : "id"`
    Name    string	    `bson: "name" json : "name"`
    Email   string	    `bson: "email" json : "email"`
}
