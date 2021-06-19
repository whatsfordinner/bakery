package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/whatsfordinner/bakery/pkg/orders"
)

type newOrder struct {
	Pastry   string `json:"pastry"`
	Customer string `json:"customer"`
}

type acceptedOrder struct {
	OrderKey string `json:"orderKey"`
}

func (a *app) homeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, "{\"message\":\"reception is attended\"}")
}

func (a *app) newOrderHandler(w http.ResponseWriter, r *http.Request) {
	orderBytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	orderData := new(newOrder)
	err = json.Unmarshal(orderBytes, orderData)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if orderData.Customer == "" || orderData.Pastry == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("New order received. Customer: %s, Pastry: %s", orderData.Customer, orderData.Pastry)

	order := orders.NewOrder(orderData.Customer, orderData.Pastry)
	key, err := a.DB.CreateOrder(r.Context(), order)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = a.Queue.PublishOrderMessage(r.Context(), &orders.OrderMessage{OrderKey: key, Pastry: order.Pastry})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result, err := json.Marshal(acceptedOrder{key})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, string(result))
}

func (a *app) orderStatusHandler(w http.ResponseWriter, r *http.Request) {
	orderKey := mux.Vars(r)["key"]

	order, err := a.DB.ReadOrder(r.Context(), orderKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if order == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	result, err := json.Marshal(*order)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, string(result))
}
