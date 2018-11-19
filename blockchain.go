package core

import (
	"fmt"
	"github.com/dgraph-io/badger"
	"os"
	"runtime"
)

const (
	dbPath      = "./tmp/blocks_%s"
	genesisData = "First Transaction from Genesis"
)

type Blockchain struct{
	LastHash []byte
	Database *badger.DB
}

func DBexists(path string) bool {
	if _, err := os.Stat(path + "/MANIFEST"); os.IsNotExist(err) {
		return false
	}

	return true
}

func ContinueBlockChain(address string) *Blockchain {
	if DBexists(dbPath) == false {
		fmt.Println("No existing blockchain found, create one!")
		runtime.Goexit()
	}

	var lastHash []byte

	opts := badger.DefaultOptions
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		lastHash, err = item.ValueCopy(nil)

		return err
	})
	Handle(err)

	chain := Blockchain{lastHash, db}

	return &chain
}

//往Blocks数组里面新加block的方法
func (chain *Blockchain)AddBlock(data string){
	var lastHash []byte

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		lastHash, err = item.ValueCopy(nil)

		return err
	})
	Handle(err)

	newBlock := NewBlock(data, lastHash)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.CurrentBlockHash, newBlock.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), newBlock.CurrentBlockHash)

		chain.LastHash = newBlock.CurrentBlockHash

		return err
	})
	Handle(err)
}


func NewBlockchain() *	Blockchain{

	var lastHash []byte

	opts := badger.DefaultOptions
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			fmt.Println("No existing blockchain found")
			genesis := Genesis()
			fmt.Println("Genesis proved")
			err = txn.Set(genesis.CurrentBlockHash, genesis.Serialize())
			Handle(err)
			err = txn.Set([]byte("lh"), genesis.CurrentBlockHash)

			lastHash = genesis.CurrentBlockHash

			return err
		} else {
			item, err := txn.Get([]byte("lh"))
			Handle(err)
			lastHash, err = item.ValueCopy(nil)
			return err
		}
	})

	Handle(err)

	blockchain := Blockchain{lastHash, db}
	return &blockchain
}
