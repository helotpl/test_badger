package md5cache

import "github.com/dgraph-io/badger/v3"

type Md5Cache struct {
	db *badger.DB
}

func (c *Md5Cache) Open(file string) error {
	opts := badger.DefaultOptions(file)
	opts.Logger = nil
	db, err := badger.Open(opts)
	if err != nil {
		return err
	}
	c.db = db
	return nil
}

func (c *Md5Cache) Close() {
	c.db.Close()
}

func (c *Md5Cache) Set(file, md5 string) error {
	err := c.db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(file), []byte(md5))
		return err
	})
	return err
}

func (c *Md5Cache) Get(file string) (string, error) {
	var md5 []byte
	err := c.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(file))
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			md5 = append([]byte{}, val...)
			return nil
		})
		return err
	})
	if err != nil {
		return "", err
	}
	return string(md5), nil
}

func (c *Md5Cache) GetAll() (map[string]string, error) {
	ret := make(map[string]string)
	err := c.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			err := item.Value(func(v []byte) error {
				ret[string(k)] = string(v)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return ret, err
}
