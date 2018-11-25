package blockchain //Be careful of the .go file directory and the package name
//Modify package name if you move this .go file

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"errors"
	"github.com/dgraph-io/badger"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	dbPath = "./database/nodes/"  //Path to store nodes' database files.
)

// 3 functions are provided to be used for accessing the blockchain database.
// SaveNode(nodeID string) ---Create a database for the node
// SaveBlockToNode(nodeID string, block *Block) ---Store a block to the node's database
// LoadBlocksFromNode(nodeID string) ---Load all blocks stored in the node's database

func SaveNode(nodeID string) error{
	nodeDBFilePath := dbPath+nodeID
	if nodeDBExists(nodeDBFilePath) {
		err := errors.New("Node DB File Already Exists")
		print(err.Error())
		return err
	}
	db := initDatabase(nodeDBFilePath)
	db.Close()
	return nil
}

func SaveBlockToNode(nodeID string, block *Block) error{
	nodeDBFilePath := dbPath+nodeID

	if ! nodeDBExists(nodeDBFilePath) {
		err :=errors.New("Node DB file Not Exists")
		print(err.Error())
		return err
	}

	db := initDatabase(nodeDBFilePath)
	defer db.Close()

	err := db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get(block.CurrentBlockHash); err == nil {
			return nil
		}
		blockData := block.serialize()
		err := txn.Set(block.CurrentBlockHash, blockData)
		return err
	})
	handle(err)
	return err
}

func LoadBlocksFromNode(nodeID string) ([]*Block, error) {
	nodeDBFilePath := dbPath+nodeID

	if ! nodeDBExists(nodeDBFilePath) {
		err :=errors.New("Node DB file Not Exists")
		print(err.Error())
		return nil, err
	}

	db := initDatabase(nodeDBFilePath)
	defer db.Close()

	blocks := []*Block{nil}

	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			_ = item.Key() //Hash
			blockData, err := item.ValueCopy(nil)  //BlockData
			block := deserialize(blockData)
			blocks = append([]*Block{block}, blocks...)
			if err!=nil{
				return err
			}
		}
		return nil
		})
	if err!=nil{
		handle(err)
	}
	return blocks, err
}

func nodeDBExists(nodeDBFilePath string) bool {
	if _, err := os.Stat(nodeDBFilePath + "/MANIFEST"); os.IsNotExist(err) {
		return false
	}
	return true
}

func initDatabase(nodeDBFilePath string) *badger.DB {
	_, err := os.Stat(nodeDBFilePath) //create path if not exitst
	if err == nil || os.IsNotExist(err) {
		os.MkdirAll(nodeDBFilePath, os.ModePerm)
	}

	opts := badger.DefaultOptions
	opts.Dir = nodeDBFilePath
	opts.ValueDir = nodeDBFilePath

	db, err := openDB(nodeDBFilePath, opts)
	handle(err)
	return db
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

func (b *Block) serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)

	handle(err)

	return res.Bytes()
}

func deserialize(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)

	handle(err)

	return &block
}

func handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}