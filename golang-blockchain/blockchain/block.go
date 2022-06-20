package blockchain

type BlockChain struct {
	Blocks []*Block
}

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int64
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
	prevBlock := chain.Blocks[len(chain.Blocks)-1]
	newBlockPtr := CreateBlock(data, prevBlock.Hash)
	chain.Blocks = append(chain.Blocks, newBlockPtr)
}

func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

func InitBlockChain() *BlockChain {
	return &BlockChain{[]*Block{Genesis()}}
}
