package blockchain

// import (
// 	"bytes"
// 	"encoding/gob"
// 	"log"
// )

// type Block struct {
// 	Hash     []byte
// 	Data     []byte
// 	PrevHash []byte
// 	Nonce    int64
// }

// func CreateBlock(data string, prevHash []byte) *Block {
// 	block := &Block{[]byte{}, []byte(data), prevHash, 0}

// 	proof := ProofFactory(block)
// 	hash, nonce := proof.MineBlockWithPOW()
// 	proof.Block.Hash = hash[:]
// 	proof.Block.Nonce = nonce
// 	return proof.Block
// }

// func Genesis() *Block {
// 	return CreateBlock("Genesis", []byte{})
// }

// func Serialize(block *Block) []byte {
// 	var res bytes.Buffer
// 	encoder := gob.NewEncoder(&res)
// 	err := encoder.Encode(block)
// 	HandleErr(err)
// 	return res.Bytes()
// }

// func Deserialize(serializedBlock []byte) Block {
// 	var block Block
// 	decoder := gob.NewDecoder(bytes.NewReader(serializedBlock))
// 	err := decoder.Decode(&block)
// 	HandleErr(err)
// 	return block
// }

// func HandleErr(err error) {
// 	if err != nil {
// 		log.Panic(err)
// 	}
// }
