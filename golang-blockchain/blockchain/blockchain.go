package blockchain

import (
	"github.com/dgraph-io/badger"
)

const (
	dbPath = "./tmp/blocks"
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

// func (chain *BlockChain) AddBlock(data string) {
// 	var lastHash []byte
// 	err := chain.Database.View(func(txn *badger.Txn) error {
// 		item, err := txn.Get([]byte("lh"))
// 		lastHash, err = item.Value()
// 		return err
// 	})
// 	newBlock := CreateBlock(data, lastHash)
// 	err = chain.Database.Update(func(txn *badger.Txn) error {
// 		err := txn.Set([]byte("lastHash"), newBlock.Hash)
// 		txn.Set(newBlock.Hash, Serialize(newBlock))
// 		return err
// 	})
// 	chain.LastHash = newBlock.Hash
// }

func InitBlockChain() {
	// var lastHash []byte

	// var lastHash []byte

	opts := badger.DefaultOptions
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)
	Handle(err)
	// err := db.Update(func(txn *badger.Txn) error {
	// 	if _, err := txn.Get([]byte("lastHash")); err == badger.ErrKeyNotFound {
	// 		fmt.Println("No existing blockchain found! Going to create genesis block!")
	// 		genesisBlock := Genesis()
	// 		err := txn.Set(genesisBlock.Hash, Serialize((genesisBlock)))
	// 		err := txn.Set("lastHash", genesisBlock.Hash))
	// 		lastHash = genesisBlock.Hash
	// 	} else {
	// 		item, err = txn.Get("lastHash")
	// 		lastHash = item.Value()
	// 	}
	// 	return err
	// })
	// HandleErr(err)
	// blockchain := BlockchaBlockChain{lastHash, db}
	// return &blockchain
}
