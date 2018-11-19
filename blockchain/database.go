package blockchain

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/dgraph-io/badger"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	dbPath = "./tmp/blocks_%s"
)

// 2 functions are provided to be used for accessing the blockchain database.
// SaveBlock()
// LoadChain()

func dbExists() bool {
	if _, err := os.Stat(dbPath + "/MANIFEST"); os.IsNotExist(err) {
		return false
	}

	return true
}

func createPathIfNotExists() {
	_, err := os.Stat(dbPath)
	if err == nil || os.IsNotExist(err) {
		os.MkdirAll(dbPath, os.ModePerm)
	}
}

func initDatabase() *badger.DB {
	createPathIfNotExists()

	opts := badger.DefaultOptions
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	db, err := openDB(dbPath, opts)
	Handle(err)
	return db
}

func SaveBlock(block *Block) {
	db := initDatabase()
	defer db.Close()

	err := db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get(block.CurrentBlockHash); err == nil {
			return nil
		}
		blockData := block.Serialize()
		err := txn.Set(block.CurrentBlockHash, blockData)
		Handle(err)
		err = txn.Set([]byte("lh"), block.CurrentBlockHash)
		Handle(err)
		return nil
	})
	Handle(err)
}

func LoadChain() (*Blockchain, error) {
	if dbExists() {

		db := initDatabase()
		defer db.Close()

		chain := &Blockchain{[]*Block{}}
		var lastHash []byte

		err := db.View(func(txn *badger.Txn) error {
			if item, err := txn.Get([]byte("lh")); err != nil {
				return err
			} else {
				lastHash, _ = item.ValueCopy(nil)
				Hash := lastHash
				for {
					item, _ = txn.Get(Hash)
					blockData, _ := item.ValueCopy(nil)
					block := Deserialize(blockData)
					chain.Blocks = append([]*Block{block}, chain.Blocks...)
					if Hash = block.PrevBlockHash; Hash== nil {
						break
					}
				}
			}
			return nil
		})
		return chain, err
	}
	return nil, nil
}

func retry(dir string, originalOpts badger.Options) (*badger.DB, error) {
	lockPath := filepath.Join(dir, "LOCK")
	if err := os.Remove(lockPath); err != nil {
		return nil, fmt.Errorf(`removing "LOCK": %s`, err)
	}
	retryOpts := originalOpts
	retryOpts.Truncate = true
	db, err := badger.Open(retryOpts)
	return db, err
}

func openDB(dir string, opts badger.Options) (*badger.DB, error) {
	if db, err := badger.Open(opts); err != nil {
		if strings.Contains(err.Error(), "LOCK") {
			if db, err := retry(dir, opts); err == nil {
				log.Println("database unlocked, value log truncated")
				return db, nil
			}
			log.Println("could not unlock database:", err)
		}
		return nil, err
	} else {
		return db, nil
	}
}

func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)

	Handle(err)

	return res.Bytes()
}

func Deserialize(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)

	Handle(err)

	return &block
}

func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}