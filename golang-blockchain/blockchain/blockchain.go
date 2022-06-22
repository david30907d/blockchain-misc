package blockchain

import (
	"fmt"
	"log"

	"github.com/dgraph-io/badger"
)

const (
	dbPath = "./tmp/blocks"
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

func (chain *BlockChain) AddBlock(data string) {
	var lastHash []byte
	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lastHash"))
		Handle(err)
		lastHash, err = item.ValueCopy(nil)
		return err
	})
	Handle(err)
	newBlock := CreateBlock(data, lastHash)
	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte("lastHash"), newBlock.Hash)
		txn.Set(newBlock.Hash, newBlock.Serialize())
		chain.LastHash = newBlock.Hash
		return err
	})
	Handle(err)
}
func (blockchain *BlockChain) Iterator() *BlockChainIterator {
	iterator := &BlockChainIterator{blockchain.LastHash, blockchain.Database}
	return iterator
}

func InitBlockChain() *BlockChain {
	var lastHash []byte
	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		log.Fatal(err)
	}

	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte("lastHash")); err == badger.ErrKeyNotFound {
			fmt.Println("No existing blockchain found! Going to create genesis block!")
			genesisBlock := Genesis()
			err = txn.Set(genesisBlock.Hash, genesisBlock.Serialize())
			Handle(err)
			err = txn.Set([]byte("lastHash"), genesisBlock.Hash)
			Handle(err)
			lastHash = genesisBlock.Hash
		} else {
			item, err := txn.Get([]byte("lastHash"))
			Handle(err)
			// Alternatively, you could also use item.ValueCopy().
			valCopy, err := item.ValueCopy(nil)
			Handle(err)
			fmt.Printf("The answer is: %x\n", valCopy)
			lastHash = valCopy
		}
		return err
	})
	Handle(err)
	blockchain := BlockChain{lastHash, db}
	return &blockchain
}

type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func (iter *BlockChainIterator) IterBackWard() *Block {
	var valCopy *Block
	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		Handle(err)
		rawBlockByte, err := item.ValueCopy(nil)
		valCopy = Deserialize(rawBlockByte)
		iter.CurrentHash = valCopy.PrevHash
		return err
	})
	Handle((err))
	return valCopy
}
