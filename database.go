package main

import (
	"errors"
	"fmt"
	"sync"
)

type database struct {
	currentBlockNo int
	blocks         map[int]Block
	transactions   map[string][]Transaction
	addressSet     map[string]struct{}

	blockMux sync.RWMutex
	transMux sync.RWMutex
	subsMux  sync.RWMutex
}

type Database interface {
	GetCurrentBlockNumber() int
	GetBlockByNumber(blockNum int) (Block, error)
	GetBlocks() ([]Block, error)
	AddBlock(b Block) error

	AddTransaction(address string, trans Transaction) error
	GetTransactionByAddress(address string) []Transaction

	AddSubscriber(address string) error
	RemoveSubscriber(address string) error
	GetAllSubscribers() []string
}

func NewDatabase() Database {
	return &database{
		currentBlockNo: 15990945,
		blocks:         make(map[int]Block, 0),
		transactions:   make(map[string][]Transaction),
		addressSet:     make(map[string]struct{}),
	}
}

func (o *database) GetCurrentBlockNumber() int {
	o.blockMux.RLock()
	defer o.blockMux.RUnlock()
	return o.currentBlockNo
}

func (o *database) GetBlockByNumber(blockNumber int) (Block, error) {
	o.blockMux.RLock()
	defer o.blockMux.RUnlock()

	block, found := o.blocks[blockNumber]
	if !found {
		return Block{}, errors.New(fmt.Sprintf("Block number not found: %v", blockNumber))
	}

	return block, nil
}

func (o *database) GetBlocks() ([]Block, error) {
	// TODO: using ordered map or db to get by offset and limit
	blocks := make([]Block, len(o.blocks))
	for _, b := range o.blocks {
		blocks = append(blocks, b)
	}
	return blocks, nil
}

func (o *database) AddBlock(block Block) error {
	o.blockMux.Lock()
	defer o.blockMux.Unlock()

	blockNum, _ := HexToInt(block.Number)
	if blockNum > o.currentBlockNo {
		o.currentBlockNo = blockNum
	}

	o.blocks[blockNum] = block
	return nil
}

func (o *database) AddTransaction(address string, trans Transaction) error {
	o.transMux.Lock()
	defer o.transMux.Unlock()

	_, found := o.transactions[address]
	if !found {
		o.transactions[address] = make([]Transaction, 0)
	}
	o.transactions[address] = append(o.transactions[address], trans)
	return nil
}

func (o *database) GetTransactionByAddress(address string) []Transaction {
	o.transMux.RLock()
	defer o.transMux.RUnlock()

	return o.transactions[address]
}

func (o *database) AddSubscriber(address string) error {
	o.subsMux.Lock()
	defer o.subsMux.Unlock()

	o.addressSet[address] = struct{}{}
	return nil
}

func (o *database) RemoveSubscriber(address string) error {
	o.subsMux.Lock()
	defer o.subsMux.Unlock()

	delete(o.addressSet, address)
	return nil
}

func (o *database) GetAllSubscribers() []string {
	o.subsMux.RLock()
	defer o.subsMux.RUnlock()

	var subs []string
	for address, _ := range o.addressSet {
		subs = append(subs, address)
	}
	return subs
}
