package main

import (
	"fmt"
	"strings"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func genId(name string) (id string) {

	id = strings.Replace(strings.ToLower(name), " ", "_", -1)
	return
}

func main() {
	var m []bson.M

	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	db := session.DB("brandflask")

	c := db.C("brandref")

	err = c.Find(nil).All(&m)
	check(err)

	for _, value := range m {
		name, ok := value["Name"].(string)
		if ok {
			fmt.Println(name, genId(name))
			update := bson.M{"$set": bson.M{"bid": genId(name)}}
			err = c.Update(value, update)
			check(err)
		}
	}

}
