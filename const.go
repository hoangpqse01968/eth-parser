package main

const (
	EthereumGatewayURL = "https://cloudflare-eth.com"

	JsonRPCVerion = "2.0"

	// Json RPC Method
	GetBlockNumMethod      = "eth_blockNumber"
	GetBlockByNumberMethod = "eth_getBlockByNumber"

	FetchingBlockDataInterval = 12 // second
	SubscriberBufferLen       = 10
	BlockBufferLen            = 50
	DefaultBufferLen          = 10
)
