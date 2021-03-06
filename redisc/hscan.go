package redisc

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"strings"
)

type Datum struct {
	Id   string // elastic document id
	Json string // elastic json document
}

func Hscan(key string, newData chan<- Datum) error {

	var (
		myid       string
		setkeyname string
		total      int
		count      int
		cursor     int64
		items      []string
	)

	c := getRedisConn()
	defer c.Close()

	for {
		values, err := redis.Values(c.Do("HSCAN", key, cursor))

		if err != nil {
			fmt.Println("hscan error on redis.Values")
		}

		values, err = redis.Scan(values, &cursor, &items)
		if err != nil {
			fmt.Println("hscan error on redis.Scan")
		}

		// a Redis Set with a unique set of Elastic Document IDs
		strary := []string{key, "set"}
		setkeyname = strings.Join(strary, "")

		for num, item := range items {
			evenodd := num % 2
			// Grab the ID
			if evenodd == 0 {
				myid = item
				_, err = c.Do("SADD", setkeyname, item)
				if err != nil {
					fmt.Println("error on SADD")
				}
			}
			if evenodd == 1 {
				// Build the struct here and put it on a channel
				mydatum := Datum{
					Id:   myid,
					Json: item,
				}
				newData <- mydatum
			}
		}
		total = total + len(items)
		count = count + 1
		if cursor == 0 {
			break
		}
	}
	fmt.Println(setkeyname, " total = ", total/2)
	return nil
}
