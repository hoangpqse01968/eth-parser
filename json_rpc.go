package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type JsonRPCWrapper interface {
	GetLatestBlockNum() (int, error)
	GetBlockByNumber(blockNum int) (*Block, error)
}

type jsonRPCWrapper struct {
}

func NewJsonRPCWrapper() JsonRPCWrapper {
	return &jsonRPCWrapper{}
}

func (o *jsonRPCWrapper) GetLatestBlockNum() (int, error) {
	resp, err := o.sendRequest(GetBlockNumMethod)
	if err != nil {
		return 0, err
	}

	hexStr, err := json.Marshal(&resp.Result)
	if err != nil {
		return 0, err
	}

	blockNum, _ := HexToInt(string(hexStr))
	return blockNum, nil
}

func (o *jsonRPCWrapper) GetBlockByNumber(blockNum int) (*Block, error) {
	hexStr := IntToHex(blockNum)
	resp, err := o.sendRequest(GetBlockByNumberMethod, hexStr, true)
	if err != nil {
		return nil, err
	}

	var block Block
	err = json.Unmarshal(resp.Result, &block)
	if err != nil {
		return nil, err
	}

	return &block, nil
}

func (o *jsonRPCWrapper) sendRequest(method string, params ...interface{}) (*Result, error) {
	reqData, err := json.Marshal(Request{
		Id:      1,
		Jsonrpc: JsonRPCVerion,
		Method:  method,
		Params:  params,
	})
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(EthereumGatewayURL, "application/json", bytes.NewBuffer(reqData))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Get block by number failed: %v", err))
	}
	defer resp.Body.Close()

	respRaw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var respData Result
	err = json.Unmarshal(respRaw, &respData)
	if err != nil {
		return nil, err
	}

	return &respData, nil
}
