package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/kurrik/json"
	"github.com/kurrik/oauth1a"
	"github.com/kurrik/twittergo"
	"gopkg.in/mgo.v2"
	"gopkg.in/yaml.v2"
)

func readYaml() (twitter map[interface{}]interface{}, err error) {
	credentials, err := ioutil.ReadFile("account.yaml")
	if err != nil {
		return
	}

	m := make(map[interface{}]interface{})

	err = yaml.Unmarshal([]byte(credentials), &m)

	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	twitter = m["twitter"].(map[interface{}]interface{})

	//fmt.Printf("%v\n", twitter)
	fmt.Printf("[%v]\n", twitter["consumerKey"])
	fmt.Printf("[%v]\n", twitter["consumerSecret"])
	fmt.Printf("[%v]\n", twitter["accessToken"])
	fmt.Printf("[%v]\n", twitter["accessTokenSecret"])

	return
}

func LoadCredentials() (client *twittergo.Client, err error) {

	twitter, err := readYaml()

	if err != nil {
		return
	}

	config := &oauth1a.ClientConfig{
		ConsumerKey:    twitter["consumerKey"].(string),
		ConsumerSecret: twitter["consumerSecret"].(string),
	}

	user := oauth1a.NewAuthorizedConfig(twitter["accessToken"].(string), twitter["accessTokenSecret"].(string))
	client = twittergo.NewClient(config, user)

	return
}

type streamConn struct {
	client *http.Client
	resp   *http.Response
	url    *url.URL
	stale  bool
	closed bool
	mu     sync.Mutex
	// wait time before trying to reconnect, this will be
	// exponentially moved up until reaching maxWait, when
	// it will exit
	wait    int
	maxWait int
	connect func() (*http.Response, error)
}

func NewStreamConn(max int) streamConn {
	return streamConn{wait: 1, maxWait: max}
}

func (conn *streamConn) Close() {
	// Just mark the connection as stale, and let the connect() handler close after a read
	conn.mu.Lock()
	defer conn.mu.Unlock()
	conn.stale = true
	conn.closed = true
	if conn.resp != nil {
		conn.resp.Body.Close()
	}
}

func (conn *streamConn) isStale() bool {
	conn.mu.Lock()
	r := conn.stale
	conn.mu.Unlock()
	return r
}

func readStream(client *twittergo.Client, sc streamConn, path string, query url.Values,
	resp *twittergo.APIResponse, handler func([]byte), done chan bool) {

	var reader *bufio.Reader
	reader = bufio.NewReader(resp.Body)

	for {
		//we've been closed
		if sc.isStale() {
			sc.Close()
			fmt.Println("Connection closed, shutting down ")
			break
		}

		line, err := reader.ReadBytes('\n')

		if err != nil {
			if sc.isStale() {
				fmt.Println("conn stale, continue")
				continue
			}

			time.Sleep(time.Second * time.Duration(sc.wait))
			//try reconnecting, but exponentially back off until MaxWait is reached then exit?
			resp, err := Connect(client, path, query)
			if err != nil || resp == nil {
				fmt.Println(" Could not reconnect to source? sleeping and will retry ")
				if sc.wait < sc.maxWait {
					sc.wait = sc.wait * 2
				} else {
					fmt.Println("exiting, max wait reached")
					done <- true
					return
				}
				continue
			}
			if resp.StatusCode != 200 {
				fmt.Printf("resp.StatusCode = %d", resp.StatusCode)
				if sc.wait < sc.maxWait {
					sc.wait = sc.wait * 2
				}
				continue
			}

			reader = bufio.NewReader(resp.Body)
			continue
		} else if sc.wait != 1 {
			sc.wait = 1
		}
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		handler(line)
	}
}

func Connect(client *twittergo.Client, path string, query url.Values) (resp *twittergo.APIResponse, err error) {
	var (
		req *http.Request
	)
	url := fmt.Sprintf("https://stream.twitter.com%v?%v", path, query.Encode())
	fmt.Println(url)
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		err = fmt.Errorf("Could not parse request: %v\n", err)
		return
	}
	resp, err = client.SendRequest(req)
	if err != nil {
		err = fmt.Errorf("Could not send request: %v\n", err)
		return
	}
	return
}

func tableExists(name string, tables []string) bool {
	for _, t := range tables {
		if name == t {
			return true
		}
	}
	return false
}

func filterStream(client *twittergo.Client, path string, dbname string, query url.Values) (err error) {
	var (
		resp *twittergo.APIResponse
	)

	sc := NewStreamConn(300)
	resp, err = Connect(client, path, query)

	done := make(chan bool)
	stream := make(chan []byte, 1000)

	go func() {
		session, err := mgo.Dial("localhost")
		if err != nil {
			panic(err)
		}
		defer session.Close()
		session.SetMode(mgo.Monotonic, true)
		db := session.DB("twitter")
		colls, err := db.CollectionNames()

		if err != nil {
			panic(err)
		}
		c := db.C(dbname)

		if !tableExists(dbname, colls) {
			info := &mgo.CollectionInfo{}
			c.Create(info)
		}

		for data := range stream {
			tweet := &twittergo.Tweet{}
			fmt.Printf("%s\n", data)
			err := json.Unmarshal(data, tweet)
			if err == nil {
				fmt.Printf("%s\n", tweet.Text())
				c.Insert(tweet)
			}
		}
	}()

	readStream(client, sc, path, query, resp, func(line []byte) {
		stream <- line
	}, done)

	return
}

func main() {
	var (
		err    error
		client *twittergo.Client
	)

	if client, err = LoadCredentials(); err != nil {
		fmt.Printf("Could not parse CREDENTIALS file: %v\n", err)
		os.Exit(1)
	}

	if len(os.Args) < 3 {
		fmt.Printf("need args\n")
		os.Exit(1)
	}

	query := url.Values{}

	for i := 2; i < len(os.Args); i++ {
		query.Add("track", os.Args[i])
	}
	fmt.Println(query)

	if err = filterStream(client, "/1.1/statuses/filter.json", os.Args[1], query); err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Printf("\n\n")
}
