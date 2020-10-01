package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/mediocregopher/radix/v3"
	"github.com/whatsfordinner/bakery/pkg/orders"
)

func TestHomeHandler(t *testing.T) {
	app, tearDown := setUp()
	defer tearDown()

	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatalf("Unable to generate HTTP request: %s", err.Error())
	}

	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected response code %d but got %d", http.StatusOK, rr.Code)
	}

	expected := `{"message":"reception is attended"}`

	if rr.Body.String() != expected {
		t.Fatalf("Message body does not match.\nExpected: %s\nGot: %s", expected, rr.Body.String())
	}
}

func TestOrderStatusHandler(t *testing.T) {
	tests := map[string]struct {
		status   int
		orderKey string
		expected *orders.Order
	}{
		"order exists": {
			http.StatusOK,
			"order0",
			&orders.Order{"cookie", "dina", "time1", "pending"},
		},
		"non-existent order": {
			http.StatusNotFound,
			"fakeorder",
			nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			app, tearDown := setUp()
			defer tearDown()

			reqURL := fmt.Sprintf("/orders/%s", test.orderKey)
			req, err := http.NewRequest("GET", reqURL, nil)

			if err != nil {
				t.Fatalf("Unable to generate HTTP request: %s", err.Error())
			}

			rr := httptest.NewRecorder()
			app.Router.ServeHTTP(rr, req)

			if rr.Code != test.status {
				t.Fatalf("Expected response code %d but got %d", test.status, rr.Code)
			}

			if test.expected != nil {
				result := new(orders.Order)
				err = json.Unmarshal(rr.Body.Bytes(), result)

				if err != nil {
					t.Fatalf("Error unmarshalling response: %s", err.Error())
				}

				if !reflect.DeepEqual(*result, *test.expected) {
					t.Fatalf("Returned objects don't match.\nExpected: %+v\nGot: %+v", *test.expected, *result)
				}
			}
		})
	}
}

func TestNewOrderHandler(t *testing.T) {
	tests := map[string]struct {
		data     string
		status   int
		expected *orders.Order
	}{
		"valid input": {
			`{"pastry":"pretzel","customer":"biggles"}`,
			http.StatusAccepted,
			&orders.Order{"pretzel", "biggles", "testtime", "pending"},
		},
		"invalid input": {
			`{"not valid":"it really isn't"}`,
			http.StatusBadRequest,
			nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			app, tearDown := setUp()
			defer tearDown()

			reqURL := "/orders"
			req, err := http.NewRequest("POST", reqURL, strings.NewReader(test.data))

			if err != nil {
				t.Fatalf("Unable to generate HTTP request: %s", err.Error())
			}

			rr := httptest.NewRecorder()
			app.Router.ServeHTTP(rr, req)

			if rr.Code != test.status {
				t.Fatalf("Expected response code %d but got %d", test.status, rr.Code)
			}

			if test.expected != nil {
				result := new(acceptedOrder)
				err := json.Unmarshal(rr.Body.Bytes(), result)

				if err != nil {
					t.Fatalf("Error unmarshalling response: %s", err.Error())
				}

				reqURL = fmt.Sprintf("/orders/%s", result.OrderKey)
				req, err = http.NewRequest("GET", reqURL, nil)

				if err != nil {
					t.Fatalf("Unable to generate HTTP request: %s", err.Error())
				}

				rr = httptest.NewRecorder()
				app.Router.ServeHTTP(rr, req)

				if rr.Code != http.StatusOK {
					t.Fatalf("Unable to find newly created order\nExpected: %d\nGot: %d", http.StatusOK, rr.Code)
				}

				order := new(orders.Order)
				err = json.Unmarshal(rr.Body.Bytes(), order)

				if err != nil {
					t.Fatalf("Error unmarshalling response: %s", err.Error())
				}

				order.OrderTime = "testtime"

				if !reflect.DeepEqual(*order, *test.expected) {
					t.Fatalf("Returned objects don't match.\nExpected: %+v\nGot: %+v", *test.expected, *order)
				}
			}
		})
	}
}

func setUp() (*app, func()) {
	// Set up the test app
	testDB := "127.0.0.1:6379"
	app := new(app)
	c := &config{
		DBHost: &testDB,
	}
	app.init(c)

	// Set up the DB
	orders := []*orders.Order{
		{"cookie", "dina", "time1", "pending"},
		{"brownie", "claude", "time2", "complete"},
		{"panini", "omar", "time3", "pending"},
	}
	db, err := radix.NewPool("tcp", testDB, 1)

	if err != nil {
		panic(err)
	}

	for i, order := range orders {
		err = db.Do(radix.FlatCmd(nil, "HSET", fmt.Sprintf("order%d", i), *order))

		if err != nil {
			panic(err)
		}
	}

	return app, func() {
		// Destroy the app
		err := app.DB.Disconnect()

		if err != nil {
			panic(err)
		}

		// Teardown the DB
		defer db.Close()
		keys := []string{}
		err = db.Do(radix.Cmd(&keys, "KEYS", "*"))
		for _, key := range keys {
			err = db.Do(radix.Cmd(nil, "DEL", key))

			if err != nil {
				panic(err)
			}
		}

	}
}
