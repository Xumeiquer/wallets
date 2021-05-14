package persistence

import (
	"encoding/json"
	"errors"

	"github.com/Xumeiquer/wallets/models"
	"go.etcd.io/bbolt"
)

// SetWallet keeps a wallet into database
func (db *Database) SetWallet(wallet models.Wallet) error {
	return db.bolt.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(MainBucket))

		data, err := json.Marshal(wallet)
		if err != nil {
			return err
		}

		// check whether the wallet is already in the DB
		w, err := db.findWalletByName(wallet.Name)
		if err != nil {
			return errors.New("wallet not added due to: " + err.Error())
		}

		if w == nil {
			// Insert the wallet in the DB
			return b.Put([]byte(wallet.ID), data)
		}

		return errors.New("wallet already added")
	})
}

func (db *Database) GetWallet(key string) (wlts []models.Wallet, err error) {
	if key == "" {
		// Get all wallets, there is no key to filter by
		err = db.bolt.View(func(tx *bbolt.Tx) error {
			b := tx.Bucket([]byte(MainBucket))
			return b.ForEach(func(k, v []byte) error {
				var w models.Wallet
				err := json.Unmarshal(v, &w)
				if err != nil {
					return err
				}
				wlts = append(wlts, w)
				return nil
			})
		})
		return
	}
	// Get wallet with ID = key
	err = db.bolt.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(MainBucket))
		value := b.Get([]byte(key))
		if value != nil {
			var w models.Wallet
			err := json.Unmarshal(value, &w)
			if err != nil {
				return err
			} else {
				wlts = append(wlts, w)
				return nil
			}
		}
		return errors.New("wallet not found")
	})
	return
}

func (db *Database) UpdateWallet(walletID string, newWallet models.Wallet) error {
	return db.bolt.Update(func(tx *bbolt.Tx) (err error) {
		b := tx.Bucket([]byte(MainBucket))
		var data []byte

		data, err = json.Marshal(&newWallet)
		if err != nil {
			return err
		}

		return b.Put([]byte(walletID), data)
	})
}

func (db *Database) DeleteWallet(walletID string) error {
	return db.bolt.Update(func(tx *bbolt.Tx) (err error) {
		b := tx.Bucket([]byte(MainBucket))
		value := b.Get([]byte(walletID))
		if value != nil {
			return b.Delete([]byte(walletID))
		}
		return errors.New("wallet not found")
	})
}

func (db *Database) findWalletByName(walletName string) (wallet *models.Wallet, err error) {
	wallet = nil
	err = db.bolt.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(MainBucket))
		c := b.Cursor()

		// loop through the bucket and check whether the wallet is there
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var w models.Wallet
			err := json.Unmarshal(v, &w)
			if err != nil {
				return err
			}
			if w.Name == walletName {
				wallet = &w
				break
			}
		}

		return nil
	})
	return
}
