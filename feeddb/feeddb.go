package main

import (
	"fmt"
	"os"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"encoding/json"
)

func main() {
	
	var start,nrow int
	start = 0
	nrow = 10
	
	var err error
	
	args:= os.Args
	if len(args) > 2 {
		start,err = strconv.Atoi(args[1])
		if err != nil {
			panic(err)
		}
		
		nrow,err = strconv.Atoi(args[2])
		if err != nil {
			panic(err)
		}
	}
	
	if len(args) > 1 {
		start,err = strconv.Atoi(args[1])
		if err != nil {
			panic(err)
		}
	}
	
	session, err := mgo.Dial("brandflask.com")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	collection := session.DB("brandflask").C("feeddb")

	records := []bson.M{}

	query := collection.Find(nil)
	
	if start > 0 {
			query = query.Skip(start)
	}
	
//	err = collection.Find(nil).All(&records).Limit(10)
	iter := query.Limit(nrow)
	
	err = iter.All(&records)
	if err != nil {
		panic(err)
	}

	result,_ := json.Marshal(records)
	fmt.Println(string(result))
	/*
	for _, row := range records {
		fmt.Printf("%v\n", row)
	}
	*/
}
