package blockchain

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/dgraph-io/badger"
)

const (
	dbFile = "/tmp/badger/MANIFEST"
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

func ContinueBlockChain(address string) *BlockChain {
	if !DBexists() {
		fmt.Println("No existing blockchain found, create one!")
		runtime.Goexit()
	}

	var lastHash []byte

	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lastHash"))
		Handle(err)
		// Alternatively, you could also use item.ValueCopy().
		lastHash, err = item.ValueCopy(nil)
		return err
	})
	Handle(err)

	chain := BlockChain{lastHash, db}

	return &chain
}

func (chain *BlockChain) AddBlock(transactions []*Transaction) {
	var lastHash []byte
	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lastHash"))
		Handle(err)
		lastHash, err = item.ValueCopy(nil)
		return err
	})
	Handle(err)
	newBlock := CreateBlock(transactions, lastHash)
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

func InitBlockChain(address string) *BlockChain {
	var lastHash []byte
	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		log.Fatal(err)
	}

	err = db.Update(func(txn *badger.Txn) error {
		fmt.Println("No existing blockchain found! Going to create genesis block!")
		genesisBlock := Genesis(address)
		err = txn.Set(genesisBlock.Hash, genesisBlock.Serialize())
		Handle(err)
		err = txn.Set([]byte("lastHash"), genesisBlock.Hash)
		Handle(err)
		lastHash = genesisBlock.Hash
		return err
	})
	Handle(err)
	blockchain := BlockChain{lastHash, db}
	return &blockchain
}

func (chain *BlockChain) ShowBalance(account_address string) int {
	transactions := chain.FindUnspentTransactions(account_address)
	balance := 0
	for _, trx := range transactions {
		for _, output := range trx.Outputs {
			if output.CanBeUnlocked(account_address) {
				balance += output.Value
			}
		}
	}
	return balance
}

func (chain *BlockChain) FindSpendableOutputs(account_address string, amount int) (int, map[string][]int) {
	result := make(map[string][]int)
	fee := 0
	unspentTrxs := chain.FindUnspentTransactions(account_address)
	for _, trx := range unspentTrxs {
		txIdStr := hex.EncodeToString(trx.ID)
		for outputIdx, output := range trx.Outputs {
			result[txIdStr] = append(result[txIdStr], outputIdx)
			fee += output.Value
			if fee >= amount {
				break
			}
		}
	}
	return fee, result
}

func (chain *BlockChain) FindUnspentTransactions(account_address string) []Transaction {
	// iterate through blockchain backward, so that we can get those trasactions haven't been spent!
	// the main logic is as follow:
	// 1. iter through blockchain backward
	// 2. iter through all the transactions in each block
	// 3. iterate through both output field and input field
	// 3.1 On output field side, it stands for your balance
	// 3.2 On input field side, it stands for how many money you've spent
	// So the main logic is as follow: iterate through input and maitain a Map, and then also iterate through Output field. If they've appeared in that Map it means you've spent it already. The rest you be unspent output, which is your balance!
	var unspentTransactions []Transaction
	spentTxOutputMap := make(map[string][]int)
	iter := chain.Iterator()
	for {
		block := iter.IterBackWard()
		// iter through the tx
		for _, transaction := range block.Transactions {
			// iter throught each transaction's input and output
			txIdStr := hex.EncodeToString(transaction.ID)
		Outputs:
			for outIdx, output := range transaction.Outputs {
				if spentTxOutputMap[txIdStr] != nil {
					for _, transactionIdx := range spentTxOutputMap[txIdStr] {
						if transactionIdx == outIdx {
							// it means this TxOutput has been spent, we shouldn't count them twice, otherwise the balance would be wrong
							continue Outputs
						}
					}
				}
				if output.CanBeUnlocked(account_address) == true {
					unspentTransactions = append(unspentTransactions, *transaction)
				}
			}
			// Skip coinbase transaction since it's simply sending rewards to miner, has nth to do with "real" transaction
			if transaction.IsCoinbase() == false {
				// first, filter out those outputs referenced by input. In other words, they've already been spent!
				for _, input := range transaction.Inputs {
					if input.CanUnlock(account_address) {
						trxIdStr := hex.EncodeToString(input.ID)
						spentTxOutputMap[trxIdStr] = append(spentTxOutputMap[trxIdStr], input.OutIdx)
					}
				}
			}
		}
		if len(block.PrevHash) == 0 {
			break
		}
	}
	return unspentTransactions
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

func DBexists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}
