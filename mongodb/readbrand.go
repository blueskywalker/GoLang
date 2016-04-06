
package main;

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
    "fmt"
    "strconv"
    "os"
     _ "encoding/json"
)

func main() {

    if len(os.Args) < 2 {
        fmt.Println("need bid")
        os.Exit(1)
    }

    var bid int64

    bid,_ = strconv.ParseInt(os.Args[1],10 , 64)

	session,err := mgo.Dial("brandflask.com")
	if err != nil {
		panic(err)
	}
	defer session.Clone()

    collection := session.DB("brandflask").C("feeddb")

    fmt.Println(bid)
    //query:= collection.Find(bson.M{"brandid":3544643373190810600})
    query:= collection.Find(bson.M{"brandid":bid})
    count, _ := query.Count()
    fmt.Println(count)

/*
    results := []bson.M{}
    err = collection.Find(bson.M{"brandid":bid}).Limit(2).All(&results)
    //err = collection.Find(bson.M{"brandid":3544643373190810456}).Limit(2).All(&results)
    //err = collection.Find(bson.M{"source":"bloomingdale"}).Limit(2).All(&results)
    if err != nil { panic(err) }

    output, _ := json.Marshal(results)
    fmt.Println(string(output))

*/
}
