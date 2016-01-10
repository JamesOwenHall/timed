package main

import (
	"fmt"
	"time"

	"github.com/JamesOwenHall/timed/cassandra"
	"github.com/gocql/gocql"
)

func main() {
	cluster := gocql.NewCluster("192.168.99.100")
	cluster.Keyspace = "timed"
	session, err := cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	defer session.Close()
	fmt.Println("Connected yo!")

	source := cassandra.Source{
		Session:     session,
		Consistency: gocql.One,
		Name:        "data",
		TimeKey:     "period",
	}

	start := time.Date(2015, 1, 2, 0, 0, 0, 0, time.UTC)
	end := time.Date(2015, 1, 5, 0, 0, 0, 0, time.UTC)
	iter := source.Iterator(start, end)

	for {
		rec, err := iter.Next()
		if err != nil {
			fmt.Println("Error:", err.Error())
		}

		if rec == nil {
			break
		}

		fmt.Println("Record:", rec)
	}
}
