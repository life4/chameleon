package chameleon

import (
	"fmt"
	"time"

	"go.etcd.io/bbolt"
)

type Database struct {
	db *bbolt.DB
}

func (db *Database) Open(path string) error {
	var err error
	opts := &bbolt.Options{Timeout: 5 * time.Second}
	db.db, err = bbolt.Open(path, 0600, opts)
	if err != nil {
		return fmt.Errorf("cannot open db: %v", err)
	}

	err = db.db.Update(func(tx *bbolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte("views"))
		if err != nil {
			return fmt.Errorf("cannot create bucket: %v", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("cannot execute transaction: %v", err)
	}

	return nil
}

func (db Database) Close() error {
	if db.db == nil {
		return nil
	}
	return db.db.Close()
}

func (db Database) Views(path Path) *Views {
	if db.db == nil {
		return nil
	}
	return &Views{
		db:   db.db,
		path: path,
	}
}
