package main

import "log"

type Subscriber struct {
	address string
	blockCh chan Block
	db      Database
}

func NewSubscriber(addr string, db Database) *Subscriber {
	s := &Subscriber{
		address: addr,
		blockCh: make(chan Block, SubscriberBufferLen),
		db:      db,
	}

	go s.run()
	go s.loadHistoryTrans()

	return s
}

func (o *Subscriber) run() {
	for block := range o.blockCh {
		for _, trans := range block.Transactions {
			if trans.From == o.address || trans.To == o.address {
				o.db.AddTransaction(o.address, trans)
			}
		}
	}
}

func (o *Subscriber) loadHistoryTrans() {
	blocks, err := o.db.GetBlocks()
	if err != nil {
		log.Printf("Load history transaction of %v failed err=%v", o.address, err)
	}
	
	for _, block := range blocks {
		o.blockCh <- block
	}
}

func (o *Subscriber) Notify(block Block) {
	o.blockCh <- block
}

func (o *Subscriber) Close() {
	close(o.blockCh)
}
