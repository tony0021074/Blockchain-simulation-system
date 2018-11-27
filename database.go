package main
//Be careful of the .go file directory and the package name
//Modify package name if you move this .go file

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/dgraph-io/badger"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// 2 functions are provided to be used for accessing the blockchain database.
// SaveBlock(block *Block, nodeID string) error ---Create a database and store, or store a block to the existing database
// LoadChain(nodeID string) ([]*Block, error) ---Load all blocks (ordering depends whether the chain is complete) stored in the node's database

const (
	dbPath = "./database/nodes/"  //Path to store nodes' database files.
)

//Logger to reset badgerdb's loggger
type nopLog struct {
	*log.Logger
}

func SaveBlock(block *Block, nodeID string) error{
	nodeDBFilePath := dbPath+nodeID
	db := initDatabase(nodeDBFilePath)
	defer db.Close()

	err := db.Update(func(txn *badger.Txn) error {
		//if current block found
		if _, err := txn.Get(block.CurrentBlockHash); err == nil {
			return errors.New("The Block Already Exists In Node DB File")
		}

		blockData := block.serialize()
		err := txn.Set(block.CurrentBlockHash, blockData)
		return err
	})
	if err !=nil{
		println(err.Error())
	}
	return err
}

func LoadChain(nodeID string) ([]*Block, error) {
	nodeDBFilePath := dbPath+nodeID
	if ! nodeDBExists(nodeDBFilePath) {
		err :=errors.New("Node DB file Not Exists")
		println(err.Error())
		return nil, err
	}
	db := initDatabase(nodeDBFilePath)
	defer db.Close()

	blocks := []*Block{}

	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		// Loop over all data in DB
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			_ = item.KeyCopy(nil) //hash
			blockData, err := item.ValueCopy(nil)  //blockData
			handle(err)
			block := deserialize(blockData)
			blocks = append([]*Block{block}, blocks...)
		}
		return nil
		})
	if err != nil{
		handle(err)
	}
	blocks = sortBlocks(blocks)
	return blocks, err
}

func nodeDBExists(nodeDBFilePath string) bool {
	if _, err := os.Stat(nodeDBFilePath + "/MANIFEST"); os.IsNotExist(err) {
		return false
	}
	return true
}

func initDatabase(nodeDBFilePath string) *badger.DB {
	//Create path if not exist
	_, err := os.Stat(nodeDBFilePath)
	if err == nil || os.IsNotExist(err) {
		os.MkdirAll(nodeDBFilePath, os.ModePerm)
	}

	nopLogger := &nopLog{Logger: log.New(os.Stderr, "", log.LstdFlags)} //Reset logger
	badger.SetLogger(nopLogger)

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

// Sort the blocks with bubble sort
func sortBlocks(blocks []*Block) []*Block {
	l := len(blocks)

	for  i := 0; i < l; i++{
		for j := 0; j < (l-1-i); j++{
			if !bytes.Equal(blocks[j].CurrentBlockHash, blocks[j+1].PrevBlockHash){
				blocks[j], blocks[j+1] = blocks[j+1], blocks[j]
			}
		}
	}

	return blocks
}

func (l *nopLog) Errorf(f string, v ...interface{}) {
	// noop
}

func (l *nopLog) Infof(f string, v ...interface{}) {
	// noop
}

func (l *nopLog) Warningf(f string, v ...interface{}) {
	// noop
}

func handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}