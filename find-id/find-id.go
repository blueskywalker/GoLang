package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"unicode/utf8"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func Levenshtein(a, b string) int {
	f := make([]int, utf8.RuneCountInString(b)+1)

	for j := range f {
		f[j] = j

	}

	for _, ca := range a {
		j := 1
		fj1 := f[0] // fj1 is the value of f[j - 1] in last iteration
		f[0]++
		for _, cb := range b {
			mn := min(f[j]+1, f[j-1]+1) // delete & insert
			if cb != ca {
				mn = min(mn, fj1+1) // change

			} else {
				mn = min(mn, fj1) // matched

			}

			fj1, f[j] = f[j], mn // save f[j] to fj1(j is about to increase), update f[j] to mn
			j++

		}

	}

	return f[len(f)-1]

}

func min(a, b int) int {
	if a <= b {
		return a
	} else {
		return b
	}
}
func max(a, b int) int {
	if a >= b {
		return a
	} else {
		return b
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func findmax(brands []string, brand string) (found string) {
	var minScore int

	minScore = 1000
	for i := 0; i < len(brands); i++ {
		score := Levenshtein(brands[i], brand)
		//fmt.Println(brand, brands[i], score)
		if score < minScore {
			minScore = score
			found = brands[i]
		}
	}

	return found
}

func findequal(brands []string, brand string) int {
	for index, v := range brands {
		if v == brand {
			return index
		}
	}

	return -1
}

func genId(name string) (id string) {

	id = strings.Replace(strings.ToLower(name), " ", "_", -1)
	return
}

func main() {
	var (
		m      []bson.M
		brands []string
	)

	f, err := os.Open("brands.txt")
	check(err)
	defer f.Close()

	r := bufio.NewReader(f)

	for {
		brand, err := r.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)

		}
		//		fmt.Print(brand)
		brands = append(brands, strings.TrimSpace(strings.ToLower(brand)))
	}

	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	db := session.DB("brandflask")

	c := db.C("nordstrom_feed")

	err = c.Find(nil).All(&m)

	if err != nil {
		log.Fatal(err)
	}

	for _, value := range m {
		name, ok := value["brand"].(string)
		if ok {
			found := findequal(brands, strings.ToLower(name))
			if found != -1 {
				fmt.Printf("%v\t%v\n", name, genId(brands[found]))
			} else {
				fmt.Printf("%v\tnot found\n", name)
			}
		}

	}

}
