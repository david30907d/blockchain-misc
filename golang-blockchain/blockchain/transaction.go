package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/david30907d/blockchain-misc/golang-blockchain/wallet"
)

const miningReward = 100

type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

func (tx *Transaction) Serialize() []byte {
	var buff bytes.Buffer
	e := gob.NewEncoder(&buff)
	err := e.Encode(tx)
	Handle(err)
	return buff.Bytes()
}

func (tx *Transaction) Hash() []byte {
	var hash [32]byte
	txCopy := *tx
	txCopy.ID = []byte{}
	hash = sha256.Sum256(txCopy.Serialize())
	return hash[:]
}

func NewTransaction(from, to string, amount int, UTXO *UTXOSet) *Transaction {
	// about outputs: it would at least has one Output struct. It stands for the receipt that belongs to `to` this address.
	// There's chance that we have 2 Outputs in Transaction{}, if the money you provided is greater than amount than you would get $xxx in change.
	var inputs []TxInput
	var outputs []TxOutput
	wallets, err := wallet.CreateWallets()
	w := wallets.GetWallet(from)
	fromPubKeyHash := wallet.PublicKeyHash(w.PublicKey)
	Handle(err)

	accumulateBalance, validOutputs := UTXO.FindSpendableOutputs(fromPubKeyHash, amount)
	if accumulateBalance < amount {
		log.Panic("You don't have enough money!")
	}
	inputs = spendYourOutputs(w, fromPubKeyHash, inputs, validOutputs)
	outputs = giveYourMoneyToPeople(outputs, to, amount)
	if accumulateBalance > amount {
		outputs = append(outputs, *NewTXOutput(accumulateBalance-amount, from))
	}
	tx := Transaction{nil, inputs, outputs}
	tx.ID = tx.Hash()
	UTXO.Blockchain.SignTransaction(&tx, w.PrivateKey)
	return &tx
}

func spendYourOutputs(w wallet.Wallet, fromPubKeyHash []byte, inputs []TxInput, validOutputs map[string][]int) []TxInput {
	// Usually inputs would be plural, since you might want to spent multiple small `Outputs` to buy an expensive stuff!
	for txid, outputIdxs := range validOutputs {
		txIdStr, err := hex.DecodeString(txid)
		Handle(err)
		for _, outputIdx := range outputIdxs {
			input := TxInput{txIdStr, outputIdx, nil, w.PublicKey}
			inputs = append(inputs, input)
		}
	}
	return inputs
}

func giveYourMoneyToPeople(outputs []TxOutput, to string, amount int) []TxOutput {
	outputs = append(outputs, *NewTXOutput(amount, to))
	return outputs
}

func CoinbaseTx(to, account_address string) *Transaction {
	if account_address == "" {
		randData := make([]byte, 20)
		_, err := rand.Read(randData)
		Handle(err)

		account_address = fmt.Sprintf("%x", randData)
	}
	txInput := TxInput{[]byte{}, -1, nil, []byte(account_address)}
	txOutput := NewTXOutput(miningReward, to)
	tx := Transaction{[]byte{}, []TxInput{txInput}, []TxOutput{*txOutput}}
	tx.ID = tx.Hash()
	return &tx
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].TrxID) == 0 && tx.Inputs[0].OutIdx == -1
}

func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	// Use your private key to sign the transaction, for the sake of letting other validator to know it you're qualified to spend these outputs!
	if tx.IsCoinbase() {
		return
	}

	for _, in := range tx.Inputs {
		if prevTXs[hex.EncodeToString(in.TrxID)].ID == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := tx.TrimmedCopy()

	for inId, in := range txCopy.Inputs {
		prevTX := prevTXs[hex.EncodeToString(in.TrxID)]
		txCopy.Inputs[inId].Signature = nil
		txCopy.Inputs[inId].PubKey = prevTX.Outputs[in.OutIdx].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Inputs[inId].PubKey = nil

		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.ID)
		Handle(err)
		signature := append(r.Bytes(), s.Bytes()...)

		tx.Inputs[inId].Signature = signature

	}
}

func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}

	for _, in := range tx.Inputs {
		if prevTXs[hex.EncodeToString(in.TrxID)].ID == nil {
			log.Panic("Previous transaction not correct")
		}
	}

	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for inId, in := range tx.Inputs {
		prevTx := prevTXs[hex.EncodeToString(in.TrxID)]
		txCopy.Inputs[inId].Signature = nil
		txCopy.Inputs[inId].PubKey = prevTx.Outputs[in.OutIdx].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Inputs[inId].PubKey = nil

		r := big.Int{}
		s := big.Int{}

		sigLen := len(in.Signature)
		r.SetBytes(in.Signature[:(sigLen / 2)])
		s.SetBytes(in.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(in.PubKey)
		x.SetBytes(in.PubKey[:(keyLen / 2)])
		y.SetBytes(in.PubKey[(keyLen / 2):])

		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubKey, txCopy.ID, &r, &s) == false {
			return false
		}
	}

	return true
}

func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	for _, in := range tx.Inputs {
		inputs = append(inputs, TxInput{in.TrxID, in.OutIdx, nil, nil})
	}

	for _, out := range tx.Outputs {
		outputs = append(outputs, TxOutput{out.Value, out.PubKeyHash})
	}

	txCopy := Transaction{tx.ID, inputs, outputs}

	return txCopy
}

func (tx Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction %x:", tx.ID))
	for i, input := range tx.Inputs {
		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TXID:     %x", input.TrxID))
		lines = append(lines, fmt.Sprintf("       Out:       %d", input.OutIdx))
		lines = append(lines, fmt.Sprintf("       Signature: %x", input.Signature))
		lines = append(lines, fmt.Sprintf("       PubKey:    %x", input.PubKey))
	}

	for i, output := range tx.Outputs {
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:  %d", output.Value))
		lines = append(lines, fmt.Sprintf("       Script: %x", output.PubKeyHash))
	}

	return strings.Join(lines, "\n")
}
