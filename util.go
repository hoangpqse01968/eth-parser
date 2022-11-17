package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type Block struct {
	Number       string
	Hash         string
	Timestamp    string
	Transactions []Transaction
}

type Transaction struct {
	BlockHash       string
	BlockNumber     string
	From            string
	To              string
	TransactionHash string
	Value           string
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Request struct {
	Id      int64         `json:"id"`
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type Result struct {
	Id      int64           `json:"id"`
	Jsonrpc string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	error   Error
}

func HexToInt(hexStr string) (int, error) {
	cleaned := strings.Replace(hexStr, "0x", "", -1)
	cleaned = strings.Replace(cleaned, "\"", "", -1)

	result, err := strconv.ParseInt(cleaned, 16, 64)
	if err != nil {
		return 0, err
	}

	return int(result), nil
}

func IntToHex(num int) string {
	return fmt.Sprintf("0x%x", num)
}
