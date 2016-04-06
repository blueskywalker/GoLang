package main

import (
	"fmt"
	_ "log"
	"gopkg.in/mgo.v2"
	_ "gopkg.in/mgo.v2/bson"
	"strings"
)

type Person struct {
	Name string
	Phone string
}

func main() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	db := session.DB("test")

	names,err := db.CollectionNames()

	if err == nil {
		fmt.Printf("%s\n",strings.Join(names,","))
	}
	/*
	c := session.DB("test").C("people")
	err = c.Insert(&Person{"Ale", "+55 53 8116 9639"},
		&Person{"Cla", "+55 53 8402 8510"})
	if err != nil {
		log.Fatal(err)
	}

	result := Person{}
	err = c.Find(bson.M{"name": "Ale"}).One(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Phone:", result.Phone)
        */
}
