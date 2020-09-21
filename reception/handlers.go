package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type newOrder struct {
	Pastry   string `json:"pastry"`
	Customer string `json:"customer"`
}

type existingOrder struct {
	Pastry   string `json:"pastry"`
	Customer string `json:"customer"`
	OrderID  string `json:"orderId"`
	Status   string `json:"status"`
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, "{\"message\":\"reception is attended\"}")
}

func newOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	orderBytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var order newOrder
	err = json.Unmarshal(orderBytes, &order)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("New order received. Customer: %s, Pastry: %s", order.Customer, order.Pastry)
}

func orderStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
