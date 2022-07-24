package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/dgraph-io/badger"
)

const (
	dbPath = "./tmp/blocks_%s"
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

func ContinueBlockChain(nodeId string) *BlockChain {
	path := fmt.Sprintf(dbPath, nodeId)
	if !DBexists(path) {
		fmt.Println("No existing blockchain found, create one!")
		runtime.Goexit()
	}
	opts := badger.DefaultOptions(path)
	var lastHash []byte
	// db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	db, err := openDB(path, opts)
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

func (chain *BlockChain) AddBlock(block *Block) {
	// before we pack up this block, need to verify all the transactions first!

	// In reality:
	// TL;DR: Each block only needs to be verified once. Upon receiving a new block, only the reference to the parent and the validity of the new block need to be checked.

	// When one gets started with Bitcoin, the client or mining software will download and verify the blockchain. During this synchronization, each block starting from the genesis block will be verified by the client. This is only necessary once for each block, because new blocks always reference the hash of the preceeding block.

	// I.e. when you are trying to verify Block 5, it will contain the hash of it's parent, Block 4. As your client has already verified that Block 4 is valid, and if the hash featured in the new Block 5 matches the known hash of Block 4, it can go straight ahead and only check whether the new Block 5 is valid.
	err := chain.Database.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get(block.Hash); err == nil {
			return nil
		}

		blockData := block.Serialize()
		err := txn.Set(block.Hash, blockData)
		Handle(err)

		item, err := txn.Get([]byte("lastHash"))
		Handle(err)
		lastHash, _ := item.ValueCopy(nil)

		item, err = txn.Get(lastHash)
		Handle(err)
		lastBlockData, _ := item.ValueCopy(nil)

		lastBlock := Deserialize(lastBlockData)

		if block.Height > lastBlock.Height {
			err = txn.Set([]byte("lastHash"), block.Hash)
			Handle(err)
			chain.LastHash = block.Hash
		}

		return nil
	})
	Handle(err)
}

func (chain *BlockChain) GetBestHeight() int {
	var lastBlock Block

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lastHash"))
		Handle(err)
		lastHash, _ := item.ValueCopy(nil)

		item, err = txn.Get(lastHash)
		Handle(err)
		lastBlockData, _ := item.ValueCopy(nil)

		lastBlock = *Deserialize(lastBlockData)

		return nil
	})
	Handle(err)

	return lastBlock.Height
}

func (chain *BlockChain) GetBlock(blockHash []byte) (Block, error) {
	var block Block

	err := chain.Database.View(func(txn *badger.Txn) error {
		if item, err := txn.Get(blockHash); err != nil {
			return errors.New("Block is not found")
		} else {
			blockData, _ := item.ValueCopy(nil)

			block = *Deserialize(blockData)
		}
		return nil
	})
	if err != nil {
		return block, err
	}

	return block, nil
}

func (chain *BlockChain) GetBlockHashes() [][]byte {
	var blocks [][]byte

	iter := chain.Iterator()

	for {
		block := iter.IterBackWard()

		blocks = append(blocks, block.Hash)

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return blocks
}

func (chain *BlockChain) MineBlock(transactions []*Transaction) *Block {
	var lastHash []byte
	var lastHeight int

	for _, tx := range transactions {
		if chain.VerifyTransaction(tx) != true {
			log.Panic("Invalid Transaction")
		}
	}
	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lastHash"))
		Handle(err)
		lastHash, err = item.ValueCopy(nil)

		item, err = txn.Get(lastHash)
		Handle(err)
		lastBlockData, _ := item.ValueCopy(nil)

		lastBlock := Deserialize(lastBlockData)

		lastHeight = lastBlock.Height

		return err
	})
	Handle(err)

	newBlock := CreateBlock(transactions, lastHash, lastHeight+1)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)
		err = txn.Set([]byte("lastHash"), newBlock.Hash)
		chain.LastHash = newBlock.Hash
		return err
	})
	Handle(err)
	return newBlock
}

func InitBlockChain(address, nodeId string) *BlockChain {
	path := fmt.Sprintf(dbPath, nodeId)
	if DBexists(path) {
		fmt.Println("Blockchain already exists")
		runtime.Goexit()
	}
	var lastHash []byte
	db, err := badger.Open(badger.DefaultOptions(path))
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

func (chain *BlockChain) FindUnspentTransactions() map[string]TxOutputs {
	// iterate through blockchain backward, so that we can get those trasactions haven't been spent!
	// the main logic is as follow:
	// 1. iter through blockchain backward
	// 2. iter through all the transactions in each block
	// 3. iterate through both output field and input field
	// 3.1 On output field side, it stands for your balance
	// 3.2 On input field side, it stands for how many money you've spent
	// So the main logic is as follow: iterate through input and maitain a Map, and then also iterate through Output field. If they've appeared in that Map it means you've spent it already. The rest you be unspent output, which is your balance!
	UTXO := make(map[string]TxOutputs)
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
				// outs := UTXO[txIdStr]
				// outs.Outputs = append(outs.Outputs, output)
				// UTXO[txIdStr] = outs
				txOutPuts := UTXO[txIdStr]
				txOutPuts.Outputs = append(txOutPuts.Outputs, output)
				UTXO[txIdStr] = txOutPuts
			}
			// Skip coinbase transaction since it's simply sending rewards to miner, has nth to do with "real" transaction
			if !transaction.IsCoinbase() {
				// first, filter out those outputs referenced by input. In other words, they've already been spent!
				for _, input := range transaction.Inputs {
					inTxID := hex.EncodeToString(input.TrxID)
					spentTxOutputMap[inTxID] = append(spentTxOutputMap[inTxID], input.OutIdx)
				}
			}
		}
		if len(block.PrevHash) == 0 {
			break
		}
	}
	return UTXO
}

func DBexists(path string) bool {
	if _, err := os.Stat(path + "/MANIFEST"); os.IsNotExist(err) {
		return false
	}

	return true
}

func (bc *BlockChain) SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey) {
	prevTXs := make(map[string]Transaction)

	for _, in := range tx.Inputs {
		prevTX, err := bc.FindTransaction(in.TrxID)
		Handle(err)
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	tx.Sign(privKey, prevTXs)
}

func (bc *BlockChain) VerifyTransaction(tx *Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}
	prevTXs := make(map[string]Transaction)

	for _, in := range tx.Inputs {
		prevTX, err := bc.FindTransaction(in.TrxID)
		Handle(err)
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	return tx.Verify(prevTXs)
}

func (bc *BlockChain) FindTransaction(trxId []byte) (Transaction, error) {
	iterator := bc.Iterator()
	for {
		block := iterator.IterBackWard()
		for _, trx := range block.Transactions {
			if bytes.Equal(trx.ID, trxId) {
				return *trx, nil
			}
		}
		if len(block.PrevHash) == 0 {
			break
		}
	}
	return Transaction{}, errors.New("Transaction does not exist")
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
