package blockchain

import (
	"blockchain/utils"
	"fmt"

	"github.com/dgraph-io/badger"
)

const (
	dbPath = "./db/blocks"
)

type BlockChain struct {
	LastBlockHash []byte
	DB            *badger.DB
}

var DBKeyLastBlockHash = []byte("last-block-hash")

func NewBlockChain() *BlockChain {
	var lastBlockHash []byte
	opts := badger.DefaultOptions(dbPath)
	db, err := badger.Open(opts)
	utils.HandlePanic(err)
	err = db.Update(func(txn *badger.Txn) error {
		_, err := txn.Get(DBKeyLastBlockHash)
		switch err {
		case nil:
			item, err := txn.Get(DBKeyLastBlockHash)
			if err != nil {
				return err
			}
			lastBlockHash, err = item.ValueCopy(nil)
			if err != nil {
				return err
			}
		case badger.ErrKeyNotFound:
			fmt.Println("No existing blockchain found")
			genesisBlock := CreateGenesisBlock()
			fmt.Println("Genesis block is created")
			err = txn.Set(genesisBlock.Hash, genesisBlock.Serialize())
			if err != nil {
				return err
			}

			err = txn.Set(DBKeyLastBlockHash, genesisBlock.Hash)
			if err != nil {
				return err
			}
		}

		return err
	})
	utils.HandlePanic(err)

	blockChain := &BlockChain{
		LastBlockHash: lastBlockHash,
		DB:            db,
	}
	return blockChain
}
func (chain *BlockChain) AddBlock(data string) {
	var lastBlockHash []byte
	err := chain.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get(DBKeyLastBlockHash)
		if err != nil {
			return err
		}
		lastBlockHash, err = item.ValueCopy(nil)
		utils.HandlePanic(err)

		return err
	})
	utils.HandlePanic(err)

	newBlock := NewBlock(data, lastBlockHash)
	err = chain.DB.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			return err
		}
		err = txn.Set(DBKeyLastBlockHash, newBlock.Hash)
		if err != nil {
			return err
		}

		// db operation and memory operation need to be
		// both success or both fail,
		// so this memory operation is involved in the transaction
		chain.LastBlockHash = newBlock.Hash
		return err
	})
	utils.HandlePanic(err)
}

type Iterator struct {
	CurrentHash []byte
	DB          *badger.DB
}

func (chain *BlockChain) NewIterator() *Iterator {
	return &Iterator{
		CurrentHash: chain.LastBlockHash,
		DB:          chain.DB,
	}
}

func (iter *Iterator) Next() *Block {
	var block *Block
	err := iter.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		if err != nil {
			return err
		}
		blockHash, err := item.ValueCopy(nil)
		block.Deserialize(blockHash)
		return err
	})
	utils.HandlePanic(err)

	return block
}
