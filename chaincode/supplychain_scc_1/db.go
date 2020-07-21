package main

import (
	"fmt"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
)

var mut sync.Mutex

func dbPut(key string, value string) {
	mut.Lock()
	db, _ := leveldb.OpenFile("/opt/db", nil)
	err := db.Put([]byte(key), []byte(value), nil)
	if err != nil {
		fmt.Println("Error writing to database")
	}
	db.Close()
	mut.Unlock()

}

func dbGet(key string) []byte {
	mut.Lock()
	db, _ := leveldb.OpenFile("/opt/db", nil)
	data, err := db.Get([]byte(key), nil)
	if err != nil {
		fmt.Println("Error getting from database")
	}
	db.Close()
	mut.Unlock()

	return data
}