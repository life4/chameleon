package chameleon

import (
	"encoding/binary"

	"go.etcd.io/bbolt"
)

type Views struct {
	db   *bbolt.DB
	path Path
}

func (views Views) Inc() error {
	err := views.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("views"))
		b := bucket.Get([]byte(views.path))
		var count uint32
		if len(b) > 1 {
			count = binary.BigEndian.Uint32(b)
		}
		b = make([]byte, 4)
		binary.BigEndian.PutUint32(b, count+1)
		return bucket.Put([]byte(views.path), b)
	})
	return err
}

func (views Views) Get() (uint32, error) {
	var count uint32
	err := views.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("views"))
		b := bucket.Get([]byte(views.path))
		if len(b) > 1 {
			count = binary.BigEndian.Uint32(b)
		}
		return nil
	})
	return count, err
}

func (views Views) All() (ViewStat, error) {
	result := make(ViewStat, 0)
	err := views.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("views"))
		return bucket.ForEach(func(k, v []byte) error {
			result.Add(string(k), binary.BigEndian.Uint32(v))
			return nil
		})
	})
	return result, err
}
