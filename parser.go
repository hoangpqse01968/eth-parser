package main

import (
	"log"
	"time"
)

type parser struct {
	newBlockCh      chan Block
	subscribeChan   chan string
	unsubscribeChan chan string

	client JsonRPCWrapper
	db     Database
	subs   map[string]*Subscriber
}

type Parser interface {
	// last parsed block
	GetCurrentBlock() int

	// add/remove address to observer
	Subscribe(address string) bool
	Unsubscribe(address string) bool

	// list of inbound and outbound transactions for an address
	GetTransactions(address string) []Transaction

	Close() error
}

func NewParser(db Database, client JsonRPCWrapper) Parser {
	addrList := db.GetAllSubscribers()
	subMap := make(map[string]*Subscriber, len(addrList))
	for _, addr := range addrList {
		subMap[addr] = NewSubscriber(addr, db)
	}

	p := &parser{
		newBlockCh:      make(chan Block, BlockBufferLen),
		subscribeChan:   make(chan string, DefaultBufferLen),
		unsubscribeChan: make(chan string, DefaultBufferLen),

		client: client,
		db:     db,
		subs:   subMap,
	}

	go p.fetchBlockDataCron()
	go p.run()

	return p
}

func (o *parser) getBlockData() error {
	latestBlockNum, err := o.client.GetLatestBlockNum()
	if err != nil {
		return err
	}

	currentBlockNo := o.db.GetCurrentBlockNumber()
	for blockNum := currentBlockNo; blockNum <= latestBlockNum; blockNum++ {
		block, err := o.client.GetBlockByNumber(blockNum)
		if err != nil {
			log.Printf("Get block no %v failed: err=%v", blockNum, err)
		} else {
			o.newBlockCh <- *block
		}
	}

	return nil
}

func (o *parser) fetchBlockDataCron() {
	ticker := time.NewTicker(FetchingBlockDataInterval * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := o.getBlockData()
			if err != nil {
				log.Printf("Get block data failed: err=%v", err)
			}
		}
	}
}

func (o *parser) run() {
	for {
		select {
		case block := <-o.newBlockCh:
			log.Printf("Handle block num %v", block.Number)

			o.db.AddBlock(block)
			o.broadcast(block)
		case address := <-o.subscribeChan:
			log.Printf("Add new subscriber %s", address)
			o.db.AddSubscriber(address)
			o.subs[address] = NewSubscriber(address, o.db)
		case address := <-o.unsubscribeChan:
			log.Printf("Remove subscriber %v", address)
			o.db.RemoveSubscriber(address)
			delete(o.subs, address)
		}
	}
}

func (o *parser) broadcast(block Block) {
	for _, sub := range o.subs {
		sub.Notify(block)
	}
}

func (o *parser) GetCurrentBlock() int {
	return o.db.GetCurrentBlockNumber()
}

func (o *parser) Subscribe(address string) bool {
	o.subscribeChan <- address
	return true
}

func (o *parser) Unsubscribe(address string) bool {
	o.unsubscribeChan <- address
	return true
}

func (o *parser) GetTransactions(address string) []Transaction {
	return o.db.GetTransactionByAddress(address)
}

func (o *parser) Close() error {
	close(o.newBlockCh)
	close(o.subscribeChan)
	close(o.unsubscribeChan)
	return nil
}
