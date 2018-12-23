package idbenchmark_test

import (
	"errors"
	"log"
	"testing"

	"go.etcd.io/bbolt"
)

func bboltConnect() (db *bbolt.DB, err error) {
	bucketName := []byte(idbenchmarkKey)

	db, err = bbolt.Open("bbolt.db", 0600, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName)
		if err != nil {
			log.Printf("create bbolt bucket error: %s", err)
		}
		return err
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func runBbolt(b *testing.B, db *bbolt.DB) {
	bucketName := []byte(idbenchmarkKey)

	for n := 0; n < b.N; n++ {
		err := db.Update(func(tx *bbolt.Tx) error {
			// Retrieve the users bucket.
			// This should be created when the DB is first opened.
			b := tx.Bucket(bucketName)

			// Generate ID for the user.
			// This returns an error only if the Tx is closed or not writeable.
			// That can't happen in an Update() call so I ignore the error check.
			id, _ := b.NextSequence()
			if id == 0 {
				log.Printf("id=0")
				return errors.New("id=0")
			}
			return nil
		})
		if err != nil {
			log.Printf("Sequence bbolt error: %v", err)
			break
		}
	}
}

func BenchmarkBbolt(b *testing.B) {
	db, err := bboltConnect()
	if err != nil {
		return
	}
	defer db.Close()

	b.ResetTimer()
	runBbolt(b, db)
	b.StopTimer()
}

func BenchmarkBboltParallel(b *testing.B) {
	db, err := bboltConnect()
	if err != nil {
		return
	}
	defer db.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			runBbolt(b, db)
		}
	})
	b.StopTimer()
}
