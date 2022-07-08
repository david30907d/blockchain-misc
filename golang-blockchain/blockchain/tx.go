package blockchain

import (
	"bytes"

	"github.com/david30907d/blockchain-misc/golang-blockchain/wallet"
)

type TxOutput struct {
	Value      int
	PubKeyHash []byte
}

type TxInput struct {
	// OutIdx stands for the index in Transaction's Outputs array
	TrxID     []byte
	OutIdx    int
	Signature []byte
	PubKey    []byte
}

func NewTXOutput(value int, address string) *TxOutput {
	newTxObj := &TxOutput{value, nil}
	newTxObj.Lock([]byte(address))
	return newTxObj
}

func (in *TxInput) UsesKey(pubKeyHash []byte) bool {
	inPubKeyHash := wallet.PublicKeyHash(in.PubKey)
	return bytes.Compare(pubKeyHash, inPubKeyHash) == 0
}

func (out *TxOutput) Lock(address []byte) {
	// the reason why we need to slice fullHash from 1 to length-4 is that
	// the full hash consists version, PubKey and CheckSum
	// the first digit is num
	// the last 4 num is checksum!
	fullHash := wallet.Base58Decode(address)
	pubKeyHash := fullHash[1 : len(fullHash)-4]
	out.PubKeyHash = pubKeyHash
}

func (out *TxOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}
