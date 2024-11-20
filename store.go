package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/boltdb/bolt"
)

func getOrSetKey() error {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		return errors.New("error opening database")
	}
	defer db.Close()

	// Create the bucket if it dosen't exist

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("MyBucket"))
		if err != nil {
			return errors.New("error creating database bucket")
		}
		return nil
	})
	if err != nil {
		return err
	}

	var apiKey string
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("MyBucket"))
		bytes := bucket.Get([]byte("OPEN_API_KEY"))
		apiKey = string(bytes)
		return nil
	})

	if apiKey == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Enter your OpenAI API key")
		input, err := reader.ReadString('\n')
		if err != nil {
			return errors.New("error reading user input")
		}

		apiKey = strings.TrimSpace(input)
		err = db.Update(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte("MyBucket"))
			err = bucket.Put([]byte("OPEN_API_KEY"), []byte(apiKey))
			if err != nil {
				return errors.New("error adding key to database")
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	err = os.Setenv("OPEN_API_KEY", apiKey)
	if err != nil {
		return errors.New("error setting enviornment variable")
	}

	return nil
}

func deleteKey() error {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		return errors.New("error opening database")
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("MyBucket"))
		if bucket == nil {
			return nil
		}
		err := bucket.Delete([]byte("OPEN_API_KEY"))
		if err != nil {
			return errors.New("error deleting the key")
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
