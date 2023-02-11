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

func CleanUpTestData() {
	UserTransaction = AllTransactions{
		City:         City{},
		Transactions: []Transaction{},
	}
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

	t.Cleanup(CleanUpTestData)

	type testInput struct {
		Name           string
		Transaction    Transaction
		RespStatusCode int
	}

	// add all type of test inputs
	testInputs := []testInput{
		{"TEST 1 create", Transaction{120.5, time.Now()}, 201},
		{"TEST 2 create another transaction", Transaction{120.5, time.Now()}, 201},
		{"TEST 3 Invalid input", Transaction{}, 422},
		{"TEST 4 bad time value", Transaction{Amount: 10}, 422},
		{"TEST 5 60 min older time", Transaction{100, time.Now().Add(time.Duration(-60) * time.Minute)}, 204},
		{"TEST 6 future time input", Transaction{100.5, time.Now().Add(time.Duration(60) * time.Minute)}, 204},
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
				t.Fail()
				t.Errorf("test api call failed err: %s", err.Error())
				return
			}

			// Check the status code is what we expect.
			if resp.Code != test.RespStatusCode {
				t.Fail()
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
			t.Fail()
			t.Errorf("test api call failed err: %s", err.Error())
			return
		}

		// Check the status code is what we expect.
		if resp.Code != http.StatusBadRequest {
			t.Fail()
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
			t.Fail()
			t.Errorf("test api call failed err: %s", err.Error())
			return
		}

		// Check the status code is what we expect.
		if resp.Code != http.StatusNoContent {
			t.Fail()
			t.Errorf("handler returned wrong status code: got %v want %v", resp.Code, http.StatusNoContent)
			return
		}
	})

}

func TestSetUserCity(t *testing.T) {

	t.Cleanup(CleanUpTestData)

	type testInput struct {
		Name           string
		City           string
		RespStatusCode int
	}

	testInputs := []testInput{
		{"TEST 1 set kochi", "kochi", 201},
		{"TEST 2 set bangalore", "bangalore", 201},
		{"TEST 3 set chennai with space", " chennai  ", 201},
		{"TEST 4 empty create", "", 400},
		{"TEST 5 empty space create", "   ", 400},
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
				t.Fail()
				t.Errorf(" %s - test api call failed err: %s", test.Name, err.Error())
				return
			}

			// Check the status code is what we expect.
			if resp.Code != test.RespStatusCode {
				t.Fail()
				t.Errorf("handler returned wrong status code: got %v want %v", resp.Code, test.RespStatusCode)
				return
			}

			// check response body if the status is 201
			if test.RespStatusCode == http.StatusCreated {
				expectedOut := "location updated"
				var testResponseBody string
				if err := json.NewDecoder(resp.Body).Decode(&testResponseBody); err != nil {
					t.Fail()
					t.Errorf("decode api response failed err: %s", err.Error())
					return
				}
				if testResponseBody != expectedOut {
					t.Fail()
					t.Errorf("returned wrong response body: got %+v expected %+v", testResponseBody, expectedOut)
					return
				}
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
			t.Fail()
			t.Errorf("test api call failed err: %s", err.Error())
			return
		}

		// Check the status code is what we expect.
		if resp.Code != http.StatusResetContent {
			t.Fail()
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
				t.Fail()
				t.Errorf("test api call failed err: %s", err.Error())
				return
			}

			// Check the status code is what we expect.
			if resp.Code != test.RespStatusCode {
				t.Fail()
				t.Errorf("handler returned wrong status code: got %v want %v", resp.Code, test.RespStatusCode)
				return
			}
		})
	}
}
func TestGetStaticsWithCityAndWithoutData(t *testing.T) {

	// t.Log("TestGetStatics with set city bangalore and without transactions data")

	t.Cleanup(CleanUpTestData)

	req := request{
		Url:    "/location",
		Method: http.MethodPost,
		Header: nil,
		body:   &City{"bangalore"},
	}

	_, err := req.CallHandler(SetUserCity)
	if err != nil {
		t.Fail()
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
				t.Fail()
				t.Errorf("test api call failed err: %s", err.Error())
				return
			}

			// Check the status code is what we expect.
			if resp.Code != test.RespStatusCode {
				t.Fail()
				t.Errorf("handler returned wrong status code: got %v want %v", resp.Code, test.RespStatusCode)
				return
			}
		})

	}

}
func TestGetStaticsWithoutCityAndWithData(t *testing.T) {

	t.Cleanup(CleanUpTestData)

	// insert transaction data
	transactions := []Transaction{
		{200.5, time.Now()},
		{10.5, time.Now()},
		{10000, time.Now()},
	}

	for _, transaction := range transactions {
		req := request{
			Url:    "/transaction",
			Method: http.MethodPost,
			Header: nil,
			body:   &transaction,
		}

		resp, err := req.CallHandler(CreateTransaction)
		if err != nil {
			t.Fail()
			t.Errorf("CreateTransaction Failed  err: %s", err.Error())
			return
		}

		// Check the status code is what we expect.
		if resp.Code != http.StatusCreated {
			t.Fail()
			t.Errorf("CreateTransaction handler returned wrong status code: got %v want %v", resp.Code, http.StatusCreated)
			return
		}
	}

	type testInput struct {
		Name           string
		City           string
		RespStatusCode int
	}
	// DO test on data
	testInputs := []testInput{
		{"TEST 1 with location", "bangalore", 200},
		{"TEST 2 with location", "kochi", 200},
		{"TEST 3 no location", "", 200},
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
				t.Fail()
				t.Errorf("test api call failed err: %s", err.Error())
				return
			}

			// Check the status code is what we expect.
			if resp.Code != test.RespStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", resp.Code, test.RespStatusCode)
				return
			}

			expectedOut := Statics{
				Sum:   10211,
				Avg:   3403.667,
				Max:   10000,
				Min:   10.5,
				Count: 3,
			}
			var testResponseBody Statics

			if err := json.NewDecoder(resp.Body).Decode(&testResponseBody); err != nil {
				t.Fail()
				t.Errorf("decode api response failed err: %s", err.Error())
				return
			}
			if testResponseBody != expectedOut {
				t.Fail()
				t.Errorf("returned wrong response body: got %+v expected %+v", testResponseBody, expectedOut)
				return
			}
		})
	}
}

func TestGetStaticsWithCityAndData(t *testing.T) {

	t.Cleanup(CleanUpTestData)

	// Set city data for end user
	req := request{
		Url:    "/location",
		Method: http.MethodPost,
		Header: nil,
		body:   &City{"bangalore"},
	}

	_, err := req.CallHandler(SetUserCity)
	if err != nil {
		t.Fail()
		t.Errorf("test api call failed err: %s", err.Error())
		return
	}

	// insert transaction data
	transactions := []Transaction{
		{200.5, time.Now()},
		{10.5, time.Now()},
		{10000, time.Now()},
	}

	for _, transaction := range transactions {
		req := request{
			Url:    "/transaction",
			Method: http.MethodPost,
			Header: nil,
			body:   &transaction,
		}

		resp, err := req.CallHandler(CreateTransaction)
		if err != nil {
			t.Fail()
			t.Errorf("CreateTransaction Failed  err: %s", err.Error())
			return
		}

		// Check the status code is what we expect.
		if resp.Code != http.StatusCreated {
			t.Fail()
			t.Errorf("CreateTransaction handler returned wrong status code: got %v want %v", resp.Code, http.StatusCreated)
			return
		}
	}

	type testInput struct {
		Name           string
		City           string
		RespStatusCode int
	}
	// DO test on data
	testInputs := []testInput{
		{"TEST 1 with location", "bangalore", 200},
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
				t.Fail()
				t.Errorf("test api call failed err: %s", err.Error())
				return
			}

			// Check the status code is what we expect.
			if resp.Code != test.RespStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", resp.Code, test.RespStatusCode)
				return
			}

			// Check the response body if the expected response is 200
			if test.RespStatusCode == http.StatusOK {
				expectedOut := Statics{
					Sum:   10211,
					Avg:   3403.667,
					Max:   10000,
					Min:   10.5,
					Count: 3,
				}
				var testResponseBody Statics

				if err := json.NewDecoder(resp.Body).Decode(&testResponseBody); err != nil {
					t.Fail()
					t.Errorf("decode api response failed err: %s", err.Error())
					return
				}
				if testResponseBody != expectedOut {
					t.Fail()
					t.Errorf("returned wrong response body: got %+v expected %+v", testResponseBody, expectedOut)
					return
				}
			}

		})
	}
}
