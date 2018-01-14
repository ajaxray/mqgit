// Copyright © 2018 Anis Uddin Ahmad <anis.programmer@gmail.com>

package db

import (
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

func Initiated(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

func connect(filePath string) *bolt.DB {
	connectedDb, err := bolt.Open(filePath, 0600, nil)
	if err != nil {
		fmt.Println("Failed to open database!")
		log.Fatal(err)
	}

	return connectedDb
}

// Read a string value by bucket and key
func Read(dbPath, bucket string, key []byte) string {
	db := connect(dbPath)
	defer db.Close()

	v := []byte("")
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		v = b.Get(key)
		return nil
	})

	if err != nil {
		fmt.Println("Read Failed!")
		log.Fatal(err)
	}
	return string(v)
}

// Write a string value by bucket and key
func Write(dbPath, bucket string, key, value []byte) {
	db := connect(dbPath)
	defer db.Close()

	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return fmt.Errorf("Create bucket: %s", err)
		}
		return b.Put(key, value)
	})

	if err != nil {
		fmt.Println("Write Failed!")
		log.Fatal(err)
	}
}

// Map a function for all elements of bucket
func Map(dbPath, bucket string, fn func(k, v []byte) error) {
	db := connect(dbPath)
	defer db.Close()

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		b.ForEach(fn)
		return nil
	})

	if err != nil {
		fmt.Println("Failed iterating over items of " + bucket)
		log.Fatal(err)
	}
}
