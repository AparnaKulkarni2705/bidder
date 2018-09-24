package testing

import (
    "testing"
     mgo "gopkg.in/mgo.v2"
     "github.com/bidder/models"
     "log"
     "gopkg.in/mgo.v2/bson"
     "github.com/stretchr/testify/assert"
)

var testDBSession *mgo.Session

const testDBName = "di_bidder_db_test"

func setupConnection(){
    session, err := mgo.Dial("127.0.0.1") 
    if err!=nil{
	panic(err)
    }
    testDBSession = session
}

func tearDown(){
    testDBSession.DB("di_bidder_db").C("real_estate").RemoveAll(nil)
    testDBSession.Close()
}

func TestInsertRealEstate(t *testing.T){
    defer tearDown()
    setupConnection()
    real_estate := &models.RealEstate{
	Description: "Test Description",
	InitialBid: 2222,
	ID : bson.NewObjectId(),
    }

    err := real_estate.Insert(testDBSession)
    assert.Nil(t, err)
}

func TestPlaceBid(t *testing.T){
    defer tearDown()
    setupConnection()
    real_estate_new := &models.RealEstate{
	Description: "Test Description",
	InitialBid: 2222,
	ID : bson.NewObjectId(),
    }

    err := real_estate_new.Insert(testDBSession)
    if err!=nil{
	log.Println("Error inserting real estate: ", err)
    }

    assert.Nil(t, models.PlaceBid(testDBSession, real_estate_new, "bidder1@fakemail.com", 4444))
    assert.Nil(t,models.PlaceBid(testDBSession, real_estate_new, "bidder2@fakemail.com", 4544))
    assert.Nil(t,models.PlaceBid(testDBSession, real_estate_new, "bidder3@fakemail.com", 4000))
    
    winner_email, bid_amount, err := real_estate_new.GetBidWinner(testDBSession)
    assert.Nil(t, err)
    assert.Equal(t, "bidder2@fakemail.com", winner_email)
    assert.Equal(t, 4444.01, bid_amount)
}

