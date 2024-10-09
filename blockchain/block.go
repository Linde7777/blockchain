package blockchain

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int
}

func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{
		Hash:     nil,
		Data:     []byte(data),
		PrevHash: prevHash,
		Nonce:    0,
	}
	pow := NewProof(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

type BlockChain struct {
	Blocks []*Block
}

func (chain *BlockChain) AddBlock(data string) {
	prevBlock := chain.Blocks[len(chain.Blocks)-1]
	newBlock := CreateBlock(data, prevBlock.Hash)
	chain.Blocks = append(chain.Blocks, newBlock)
}

func CreateGenesisBlock() *Block {
	return CreateBlock("Genesis Block", []byte{})
}

func InitBlockChain() *BlockChain {
	return &BlockChain{
		Blocks: []*Block{CreateGenesisBlock()},
	}
}
