// Copyright 2011 Arne Roomann-Kurrik
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/kurrik/oauth1a"
	"github.com/kurrik/twittergo"
	"gopkg.in/yaml.v2"
)

func LoadCredentials() (client *twittergo.Client, err error) {
	credentials, err := ioutil.ReadFile("account.yaml")
	if err != nil {
		fmt.Println("There is no account.yaml")
		return nil, err
	}

	m := make(map[interface{}]interface{})

	err = yaml.Unmarshal([]byte(credentials), &m)

	if err != nil {
		log.Fatalf("error: %v", err)
	}

	twitter := m["twitter"].(map[interface{}]interface{})

	//fmt.Printf("%v\n", twitter)
	fmt.Printf("[%v]\n", twitter["consumerKey"])
	fmt.Printf("[%v]\n", twitter["consumerSecret"])
	fmt.Printf("[%v]\n", twitter["accessToken"])
	fmt.Printf("[%v]\n", twitter["accessTokenSecret"])

	config := &oauth1a.ClientConfig{
		ConsumerKey:    twitter["consumerKey"].(string),
		ConsumerSecret: twitter["consumerSecret"].(string),
	}

	user := oauth1a.NewAuthorizedConfig(twitter["accessToken"].(string), twitter["accessTokenSecret"].(string))
	client = twittergo.NewClient(config, user)

	return client, err
}

func main() {
	var (
		err     error
		client  *twittergo.Client
		req     *http.Request
		resp    *twittergo.APIResponse
		results *twittergo.SearchResults
	)
	client, err = LoadCredentials()
	if err != nil {
		fmt.Printf("Could not parse account.yaml file: %v\n", err)
		os.Exit(1)
	}
	if len(os.Args) < 2 {
		fmt.Printf("need args\n")
		os.Exit(1)
	}

	query := url.Values{}

	query.Set("q", os.Args[1])

	url := fmt.Sprintf("/1.1/search/tweets.json?%v", query.Encode())
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Could not parse request: %v\n", err)
		os.Exit(1)
	}
	resp, err = client.SendRequest(req)
	if err != nil {
		fmt.Printf("Could not send request: %v\n", err)
		os.Exit(1)
	}
	results = &twittergo.SearchResults{}
	err = resp.Parse(results)
	if err != nil {
		fmt.Printf("Problem parsing response: %v\n", err)
		os.Exit(1)
	}

	for i, tweet := range results.Statuses() {
		user := tweet.User()
		fmt.Printf("%v.) %v\n", i+1, tweet.Text())
		fmt.Printf("From %v (@%v) ", user.Name(), user.ScreenName())
		fmt.Printf("at %v\n\n", tweet.CreatedAt().Format(time.RFC1123))
	}
	if resp.HasRateLimit() {
		fmt.Printf("Rate limit:           %v\n", resp.RateLimit())
		fmt.Printf("Rate limit remaining: %v\n", resp.RateLimitRemaining())
		fmt.Printf("Rate limit reset:     %v\n", resp.RateLimitReset())
	} else {
		fmt.Printf("Could not parse rate limit from response.\n")
	}
}
