package main

import (
	"flag"
	"fmt"
	badger "github.com/dgraph-io/badger/v3"
	"log"
)

func main() {
	op := flag.String("op", "set", "operation: set, get, getall")
	key := flag.String("key", "key", "key for set operation")
	value := flag.String("value", "value", "value for set operation")
	flag.Parse()

	opts := badger.DefaultOptions("database.db")
	opts.Logger = nil
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if *op == "set" {
		err := db.Update(func(txn *badger.Txn) error {
			err := txn.Set([]byte(*key), []byte(*value))
			return err
		})
		if err != nil {
			log.Fatal(err)
		}
	} else if *op == "getall" {
		err := db.View(func(txn *badger.Txn) error {
			opts := badger.DefaultIteratorOptions
			opts.PrefetchSize = 10
			it := txn.NewIterator(opts)
			defer it.Close()
			for it.Rewind(); it.Valid(); it.Next() {
				item := it.Item()
				k := item.Key()
				err := item.Value(func(v []byte) error {
					fmt.Printf("key=%s, value=%s\n", k, v)
					return nil
				})
				if err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
	} else if *op == "get" {
		err := db.View(func(txn *badger.Txn) error {
			item, err := txn.Get([]byte(*key))
			if err != nil {
				log.Fatal(err)
			}

			err = item.Value(func(val []byte) error {
				fmt.Printf("key=%s, value=%s\n", *key, val)
				return nil
			})
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}
