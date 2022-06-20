package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

const Difficulty = 18

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func ProofFactory(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))
	pow := &ProofOfWork{b, target}
	return pow
}

func GetHashOfBlock(pow *ProofOfWork) []byte {
	bytesOfData := bytes.Join([][]byte{
		pow.Block.Data, pow.Block.PrevHash, EncodeIntToHex(pow.Block.Nonce), EncodeIntToHex(Difficulty)}, []byte{})
	hash := sha256.Sum256(bytesOfData)
	return hash[:]
}

func (pow *ProofOfWork) MineBlockWithPOW() ([]byte, int64) {
	var intHash big.Int
	var hash []byte
	nonce := 0
	for nonce < math.MaxInt64 {
		hash = GetHashOfBlock(pow)
		intHash.SetBytes(hash[:])
		if intHash.Cmp(pow.Target) == -1 {
			break
		} else {
			nonce++
		}
	}
	return hash, nonce
}

func Validate(pow *ProofOfWork) bool {
	var intHash big.Int
	hash := GetHashOfBlock(pow)
	intHash.SetBytes(hash[:])
	return intHash.Cmp(pow.Target) == -1
}

func EncodeIntToHex(num int64) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, num)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(buf.Bytes())
	return buf.Bytes()
}
