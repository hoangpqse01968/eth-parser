package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func getCurrentBlockNumberHandler(p Parser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		curBlockNum := p.GetCurrentBlock()
		w.Header().Add("Content-Type", "application/json")
		resp := `{"current_block_number":"` + strconv.Itoa(curBlockNum) + `"}`
		w.Write([]byte(resp))
	}
}

func subscribeHandler(p Parser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		p.Subscribe(address)
		w.Header().Add("Content-Type", "application/json")
		resp := `{"message":"subscribed successfully"}`
		w.Write([]byte(resp))
	}
}

func unsubscribeHandler(p Parser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		p.Unsubscribe(address)
		w.Header().Add("Content-Type", "application/json")
		resp := `{"message":"unsubscribed successfully"}`
		w.Write([]byte(resp))
	}
}

func getTransactionsHandler(p Parser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		trans := p.GetTransactions(address)
		w.Header().Add("Content-Type", "application/json")
		data, err := json.Marshal(trans)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(data)
	}
}

func main() {
	db := NewDatabase()
	client := NewJsonRPCWrapper()
	p := NewParser(db, client)

	http.HandleFunc("/get-current-block-number", getCurrentBlockNumberHandler(p))
	http.HandleFunc("/subscribe", subscribeHandler(p))
	http.HandleFunc("/unsubscribe", unsubscribeHandler(p))
	http.HandleFunc("/get-transactions", getTransactionsHandler(p))

	http.ListenAndServe(":8080", nil)
}
