package persistence

import (
	"time"

	"github.com/labstack/echo"
	"go.etcd.io/bbolt"
)

const (
	// MainBucket is the bucket where the whole thing will be stored
	MainBucket = "Wallets"

	// DB is the key used for storing the datatore in the request context
	DB = "db"
)

// Database is the main datastore implementation
type Database struct {
	bolt *bbolt.DB
}

// Persistence is the middleware handler for the persistence datastore
func Persistence(db *Database) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			c.Set(DB, db)
			return next(c)
		}
	}
}

// NewDatabase creates a new database object.
func NewDatabase(path string) (*Database, error) {
	db, err := bbolt.Open(path, 0666, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	// Create the bucket just in case it does not exist
	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(MainBucket))
		if err != nil {
			return err
		}
		return nil
	})

	return &Database{
		bolt: db,
	}, err
}

// Close closes the database
func (db *Database) Close() {
	db.bolt.Close()
}

// func (db *Database) UpdateWalletFunds(walletID, fundINSI string) error {
// 	return db.bolt.Update(func(tx *bbolt.Tx) error {
// 		b := tx.Bucket([]byte(MainBucket))

// 		walletData := b.Get([]byte(walletID))
// 		if walletData != nil {
// 			var w models.Wallet
// 			err := json.Unmarshal(walletData, &w)
// 			if err != nil {
// 				return err
// 			}
// 			w.AddFund(fundINSI)
// 			walletData, err = json.Marshal(w)
// 			if err != nil {
// 				return err
// 			}
// 			return b.Put([]byte(walletID), walletData)
// 		}
// 		return errors.New("wallet not found")
// 	})
// }
