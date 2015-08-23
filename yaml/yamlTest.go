package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

func main() {

	credentials, err := ioutil.ReadFile("account.yaml")
	if err != nil {
		return
	}

	m := make(map[interface{}]interface{})

	err = yaml.Unmarshal([]byte(credentials), &m)

	if err != nil {
		log.Fatalf("error: %v", err)
	}
	twitter := m["twitter"].(map[interface{}]interface{})
	//fmt.Printf("%v\n", twitter)
	fmt.Printf("%v\n", twitter["accessTokenSecret"])

}
