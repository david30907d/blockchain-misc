package blockchain

import (
	"bytes"
	"crypto/sha256"
)

type BlockChain struct {
	blocks []*Block
}

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int64
}

func (b *Block) DeriveHash() {
	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{})
	hash := sha256.Sum256(info)
	b.Hash = hash[:]
}

func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash, 0}

	proof := ProofFactory(block)
	hash, nonce := proof.MineBlockWithPOW()
	proof.Block.Hash = hash[:]
	proof.Block.Nonce = nonce
	return proof.Block
}

func (chain *BlockChain) AddBlock(data string) {
	prevBlock := chain.blocks[len(chain.blocks)-1]
	newBlockPtr := CreateBlock(data, prevBlock.Hash)
	chain.blocks = append(chain.blocks, newBlockPtr)
}

func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

func InitBlockChain() *BlockChain {
	return &BlockChain{[]*Block{Genesis()}}
}
