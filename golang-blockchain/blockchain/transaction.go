package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

const miningReward = 100

type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

type TxOutput struct {
	Value  int
	PubKey string
}

type TxInput struct {
	// OutIdx stands for the index in Transaction's Outputs array
	ID     []byte
	OutIdx int
	Sig    string
}

func NewTransaction(from, to string, amount int, chain *BlockChain) *Transaction {
	// about outputs: it would at least has one Output struct. It stands for the receipt that belongs to `to` this address.
	// There's chance that we have 2 Outputs in Transaction{}, if the money you provided is greater than amount than you would get $xxx in change.
	var inputs []TxInput
	var outputs []TxOutput
	accumulateBalance, validOutputs := chain.FindSpendableOutputs(from, amount)
	if accumulateBalance < amount {
		log.Panic("You don't have enough money!")
	}
	inputs = spendYourOutputs(from, inputs, validOutputs)
	outputs = giveYourMoneyToPeople(outputs, to, amount)
	if accumulateBalance > amount {
		outputs = append(outputs, TxOutput{accumulateBalance - amount, from})
	}
	tx := Transaction{nil, inputs, outputs}
	tx.SetID()
	return &tx
}

func spendYourOutputs(from string, inputs []TxInput, validOutputs map[string][]int) []TxInput {
	for txid, outputIdxs := range validOutputs {
		txIdStr, err := hex.DecodeString(txid)
		Handle(err)
		for _, outputIdx := range outputIdxs {
			input := TxInput{txIdStr, outputIdx, from}
			inputs = append(inputs, input)
		}
	}
	return inputs
}

func giveYourMoneyToPeople(outputs []TxOutput, to string, amount int) []TxOutput {
	outputs = append(outputs, TxOutput{amount, to})
	return outputs
}

func CoinbaseTx(to, account_address string) *Transaction {
	if account_address == "" {
		account_address = fmt.Sprintf("Coins to %s", to)
	}
	txInput := TxInput{[]byte{}, -1, account_address}
	txOutput := TxOutput{miningReward, to}
	tx := Transaction{[]byte{}, []TxInput{txInput}, []TxOutput{txOutput}}
	tx.SetID()
	return &tx
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].OutIdx == -1
}

func (in *TxInput) CanUnlock(accountName string) bool {
	return in.Sig == accountName
}

func (out *TxOutput) CanBeUnlocked(accountName string) bool {
	return out.PubKey == accountName
}

func (tx *Transaction) SetID() {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(tx)
	Handle(err)
	hash := sha256.Sum256(buffer.Bytes())
	tx.ID = hash[:]
}
