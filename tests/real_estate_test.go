package testing

import (
    "testing"
     mgo "gopkg.in/mgo.v2"
     "github.com/bidder/models"
     "log"
     "gopkg.in/mgo.v2/bson"
     "github.com/stretchr/testify/assert"
     "sync"
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

var wg sync.WaitGroup

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

    bidders := make(map[string]float64)
    bidders["bidder1@fakemail.com"] = 4000
    bidders["bidder2@fakemail.com"] = 4001
    bidders["bidder3@fakemail.com"] = 4002
    bidders["bidder4@fakemail.com"] = 4003
    bidders["bidder5@fakemail.com"] = 4004
    bidders["bidder6@fakemail.com"] = 4005
    bidders["bidder7@fakemail.com"] = 4006
    bidders["bidder8@fakemail.com"] = 4007


    for email, bid_amount := range bidders {
	wg.Add(1)
	go  call_go_routine_to_place_bid(real_estate_new, email, bid_amount)
    }
    
    wg.Wait()

    winner_email, bid_amount, err := real_estate_new.GetBidWinner(testDBSession)
    assert.Nil(t, err)
    assert.Equal(t, "bidder8@fakemail.com", winner_email)
    assert.Equal(t, 4006.01, bid_amount)
}

func call_go_routine_to_place_bid(re *models.RealEstate, email string, bid_amount float64){
    defer wg.Done()
    models.PlaceBid(testDBSession, re, email, bid_amount)
}
