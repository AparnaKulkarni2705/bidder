package models

import (
    "log"
    "errors"
    mgo "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"

)

const REAL_ESTATE_COLLECTION_NAME = "real_estate"

type BidderBidMap struct{
    BidderEmail	string	`bson:"bidder_email" json:"bidder_email"`
    BidAmount	float64	`bson:"bid_amount" json:"bid_amount"`
}

type RealEstate struct{
    ID			bson.ObjectId   `bson:"_id,omitempty" json:"id"`
    Description		string		`bson:"description" json:"description"`
    InitialBid		float64		`bson:"initial_bid" json:"initial_bid"`
    HighestBidder	BidderBidMap	`bson:"highest_bidder" json:"highest_bidder"`
    SecondHighestBidder	BidderBidMap	`bson:"second_highest_bidder" json:"second_highest_bidder"`
    IsSold		bool		`bson:"is_sold" json:"is_sold"`
}

func (re *RealEstate) Insert(session *mgo.Session) error {
    real_estate_collection := session.DB("di_bidder_db").C(REAL_ESTATE_COLLECTION_NAME)
   
    err:=real_estate_collection.Insert(re)
    if err !=nil{
	log.Fatalf("Error inserting real_estate %s: Error: %v", re.Description, err)
	return err
    }
    return nil
}

func GetRealEstateByID(session *mgo.Session, ID bson.ObjectId) (*RealEstate, error) {
    real_estate_collection := session.DB("di_bidder_db").C(REAL_ESTATE_COLLECTION_NAME)
    var re *RealEstate

    err := real_estate_collection.Find(bson.M{"_id":ID}).One(re)
    if err!=nil{
	log.Fatalf("Error fetching real estate record with ID : %s Error: %v", ID, err)
	return nil, err
    }
    return re, nil
}

func (re *RealEstate) SetInitialBid(session *mgo.Session, initial_bid float64) error{
    real_estate_collection := session.DB("di_bidder_db").C(REAL_ESTATE_COLLECTION_NAME)
    
    filter := bson.M{"_id": re.ID}
    change := bson.M{
		"$set": bson.M{
			"initial_bid": initial_bid,
		},
	}

	err := real_estate_collection.Update(filter, change)
	if err!=nil{
	    log.Fatalf("Error updating real estate record with InitialBid Error: %v", err)
	    return err
	}
	return nil
}

func PlaceBid(session *mgo.Session, re *RealEstate, bidder_email string, bid_amount float64) error {
    real_estate_collection := session.DB("di_bidder_db").C(REAL_ESTATE_COLLECTION_NAME)
   
    real_estate_fetched := &RealEstate{}   
    err := real_estate_collection.FindId(re.ID).One(real_estate_fetched)
    if err!=nil{
	return err
    }
    filter := bson.M{"_id": re.ID}
    change := bson.M{}
    
    if bid_amount < re.InitialBid{
	log.Println("Bid amount can't be less than inital bid amount")
	return errors.New("Bid amount can't be less than inital bid amount") 
    }

    //if it is first bid, set highest and second highest bid to the initial bid
    if real_estate_fetched.HighestBidder.BidderEmail == ""{
	change = bson.M{
	    "$set":bson.M{
		"highest_bidder.bidder_email": bidder_email,
		"second_highest_bidder.bidder_email": bidder_email,
		"highest_bidder.bid_amount" : bid_amount,
		"second_highest_bidder.bid_amount": bid_amount,
	    },
	}
    } else if real_estate_fetched.HighestBidder.BidAmount < bid_amount{
	change = bson.M{
	    "$set":bson.M{
		"highest_bidder.bidder_email" : bidder_email,
		"highest_bidder.bid_amount" : bid_amount,
		"second_highest_bidder.bidder_email" : real_estate_fetched.HighestBidder.BidderEmail,
		"second_highest_bidder.bid_amount" : real_estate_fetched.HighestBidder.BidAmount,
	    },
	}
    } else if real_estate_fetched.SecondHighestBidder.BidAmount < bid_amount{
	change = bson.M{
	    "$set":bson.M{
		"second_highest_bidder.bidder_email":bidder_email,
		"second_highest_bidder.bid_amount":bid_amount,
	    },
	}
    } else{
	log.Println("No change in HighestBidder and SecondHighestBidder for real estate: ", re.ID)
	return nil
    }


    err = real_estate_collection.Update(filter, change)
    if err != nil{
	log.Fatalf("Error updating real estate record having ID: %s Error: %v", re.ID.String(), err)
	return err
    }
    return nil
}

func (re *RealEstate) GetBidWinner(session *mgo.Session) (string, float64, error){
    real_estate_collection := session.DB("di_bidder_db").C(REAL_ESTATE_COLLECTION_NAME)
    fetched_re := &RealEstate{}

    err := real_estate_collection.Find(bson.M{"_id":re.ID}).One(fetched_re)
    if err!=nil{
	log.Fatalf("Error fetching real estate record with ID : %v Error: %v", re.ID, err)
	return "", 0.0, err
    }
    //update is_sold flag
    filter := bson.M{"_id":re.ID}
    change := bson.M{
	"$set":bson.M{
	    "is_sold":true,
	},
    }

    real_estate_collection.Update(filter, change)
    winner := fetched_re.HighestBidder.BidderEmail
    amount_to_be_paid := fetched_re.SecondHighestBidder.BidAmount + 0.01
    return winner, amount_to_be_paid, nil 

}
