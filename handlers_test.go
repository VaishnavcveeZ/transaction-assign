package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type request struct {
	Url    string
	Method string
	Header map[string]string
	body   interface{}
}

func (r *request) CallHandler(handlerFunc func(w http.ResponseWriter, r *http.Request)) (*httptest.ResponseRecorder, error) {

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(r.body)
	if err != nil {
		// log.Fatal(err)
		return nil, err
	}

	req := httptest.NewRequest(r.Method, r.Url, io.NopCloser(&buf))
	if r.Header != nil {
		for k, v := range r.Header {
			req.Header.Add(k, v)
		}
	}

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(handlerFunc)

	handler.ServeHTTP(w, req)

	return w, nil
}

func TestCreateTransaction(t *testing.T) {

	type testInput struct {
		Name           string
		Transaction    Transaction
		RespStatusCode int
	}

	// add all type of test inputs
	testInputs := []testInput{
		{"TEST 1 create", Transaction{120.5, time.Now()}, 201},
		{"TEST 2 Invalid input", Transaction{}, 422},
		{"TEST 3 bad time value", Transaction{Amount: 10}, 422},
		{"TEST 4 60 min older time", Transaction{100, time.Now().Add(time.Duration(-60) * time.Minute)}, 204},
		{"TEST 5 future time input", Transaction{100.5, time.Now().Add(time.Duration(60) * time.Minute)}, 204},
	}

	// iterate test inputs and run test.
	for _, test := range testInputs {

		t.Run(test.Name, func(t *testing.T) {
			req := request{
				Url:    "/transaction",
				Method: http.MethodPost,
				Header: nil,
				body:   &test.Transaction,
			}

			resp, err := req.CallHandler(CreateTransaction)
			if err != nil {
				t.Errorf("test api call failed err: %s", err.Error())
				return
			}

			// Check the status code is what we expect.
			if resp.Code != test.RespStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", resp.Code, test.RespStatusCode)
				return
			}
		})
	}

	// This test case needs different input data
	t.Run("TEST 6 Invalid body", func(t *testing.T) {
		req := request{
			Url:    "/transaction",
			Method: http.MethodPost,
			Header: nil,
			body:   "", // pass invalid type of input
		}

		resp, err := req.CallHandler(CreateTransaction)
		if err != nil {
			t.Errorf("test api call failed err: %s", err.Error())
			return
		}

		// Check the status code is what we expect.
		if resp.Code != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", resp.Code, http.StatusBadRequest)
			return
		}
	})

}

func TestDeleteTransaction(t *testing.T) {

	t.Run("TEST 1 delete transactions", func(t *testing.T) {
		req := request{
			Url:    "/transaction",
			Method: http.MethodDelete,
			Header: nil,
			body:   nil,
		}

		resp, err := req.CallHandler(DeleteTransaction)
		if err != nil {
			t.Errorf("test api call failed err: %s", err.Error())
			return
		}

		// Check the status code is what we expect.
		if resp.Code != http.StatusNoContent {
			t.Errorf("handler returned wrong status code: got %v want %v", resp.Code, http.StatusNoContent)
			return
		}
	})

}

func TestSetUserCity(t *testing.T) {

	type testInput struct {
		Name           string
		City           string
		RespStatusCode int
	}

	testInputs := []testInput{
		{"TEST 1 set kochi", "kochi", 201},
		{"TEST 2 set bangalore", "bangalore", 201},
		{"TEST 3 empty create", "", 400},
		{"TEST 4 empty space create", "   ", 400},
	}
	for _, test := range testInputs {

		t.Run(test.Name, func(t *testing.T) {
			req := request{
				Url:    "/location",
				Method: http.MethodPost,
				Header: nil,
				body:   &City{test.City},
			}

			resp, err := req.CallHandler(SetUserCity)
			if err != nil {
				t.Errorf(" %s - test api call failed err: %s", test.Name, err.Error())
				return
			}

			// Check the status code is what we expect.
			if resp.Code != test.RespStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", resp.Code, test.RespStatusCode)
				return
			}
		})

	}
}

func TestResetUserCity(t *testing.T) {

	t.Run("TEST 1 rest city", func(t *testing.T) {
		req := request{
			Url:    "/location",
			Method: http.MethodPut,
			Header: nil,
			body:   nil,
		}

		resp, err := req.CallHandler(ResetUserCity)
		if err != nil {
			t.Errorf("test api call failed err: %s", err.Error())
			return
		}

		// Check the status code is what we expect.
		if resp.Code != http.StatusResetContent {
			t.Errorf("handler returned wrong status code: got %v want %v", resp.Code, http.StatusResetContent)
			return
		}

	})

}

func TestGetStaticsWithoutCityAndData(t *testing.T) {

	// t.Log("TestGetStatics without set city and transactions data")

	type testInput struct {
		Name           string
		City           string
		RespStatusCode int
	}

	// test inputs without set end point user city and transaction data in store
	testInputs := []testInput{
		{"TEST 1 no location", "", 204},
		{"TEST 2 with location kochi", "kochi", 204},
		{"TEST 3 with location bangalore", "bangalore", 204},
	}
	for _, test := range testInputs {

		t.Run(test.Name, func(t *testing.T) {
			req := request{
				Url:    "/statistics",
				Method: http.MethodGet,
				Header: nil,
				body:   nil,
			}

			if test.City != "" {
				req.Header = map[string]string{"location": test.City}
			}

			resp, err := req.CallHandler(GetStatics)
			if err != nil {
				t.Errorf("test api call failed err: %s", err.Error())
				return
			}

			// Check the status code is what we expect.
			if resp.Code != test.RespStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", resp.Code, test.RespStatusCode)
				return
			}
		})
	}
}
func TestGetStaticsWithCityAndWithoutData(t *testing.T) {

	// t.Log("TestGetStatics with set city bangalore and without transactions data")

	req := request{
		Url:    "/location",
		Method: http.MethodPost,
		Header: nil,
		body:   &City{"bangalore"},
	}

	_, err := req.CallHandler(SetUserCity)
	if err != nil {
		t.Errorf("test api call failed err: %s", err.Error())
		return
	}

	type testInput struct {
		Name           string
		City           string
		RespStatusCode int
	}

	// test inputs without set end point user city and transaction data in store
	testInputs := []testInput{
		{"TEST 1 with location", "bangalore", 204},
		{"TEST 2 with location", "kochi", 401},
		{"TEST 3 no location", "", 401},
	}
	for _, test := range testInputs {

		t.Run(test.Name, func(t *testing.T) {
			req := request{
				Url:    "/statistics",
				Method: http.MethodGet,
				Header: nil,
				body:   nil,
			}

			if test.City != "" {
				req.Header = map[string]string{"location": test.City}
			}

			resp, err := req.CallHandler(GetStatics)
			if err != nil {
				t.Errorf("test api call failed err: %s", err.Error())
				return
			}

			// Check the status code is what we expect.
			if resp.Code != test.RespStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", resp.Code, test.RespStatusCode)
				return
			}
		})

	}
}
