package main

import (
	"fmt"
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/olivere/elastic.v2"
)

type Brand struct {
	Oid                 string
	RawBrandName        string
	NormalizedBrandName string
	BrandID             int64
}

func main() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	db := session.DB("brandflask")

	collection := db.C("brand_normalize")

	record := []bson.M{}

	err = collection.Find(nil).All(&record)

	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	client, err := elastic.NewClient()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	for _, value := range record {
		var brand Brand

		brand.Oid = value["_id"].(bson.ObjectId).Hex()
		brand.RawBrandName = value["raw_brand_name"].(string)
		brand.NormalizedBrandName = value["normalized_brand_name"].(string)
		brand.BrandID = value["brand_id"].(int64)

		fmt.Println(brand)
	}

}
