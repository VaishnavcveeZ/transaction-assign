package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

// New struct object of AllTransactions to store data
var UserTransaction = AllTransactions{
	Transactions: []Transaction{},
}

// CreateTransaction - insert new transaction to the transactions list.
func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")

	var transaction Transaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	// check if any of the fields are not parsable.
	if transaction.Amount == 0 || transaction.TimeStamp.IsZero() {
		http.Error(w, "invalid input", http.StatusUnprocessableEntity)
		return
	}

	// load the IST location
	location, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Parse the input time stamp to ist format.
	year, month, day := transaction.TimeStamp.Date()
	hour, minute, second := transaction.TimeStamp.Clock()
	transaction.TimeStamp = time.Date(year, month, day, hour, minute, second, 0, location)

	// check if transaction date is older than 60 second or it is in the future.
	if time.Since(transaction.TimeStamp) > time.Duration(60*time.Second) || transaction.TimeStamp.After(time.Now()) {
		http.Error(w, "invalid input", http.StatusNoContent)
		return
	}

	UserTransaction.Transactions = append(UserTransaction.Transactions, transaction)

	// Pre-computing and updating of transaction data
	UserTransaction.TotalAmount += transaction.Amount
	UserTransaction.TotalTransactions += 1
	if UserTransaction.MaxAmount < transaction.Amount {
		UserTransaction.MaxAmount = transaction.Amount
	}
	if UserTransaction.TotalTransactions == 1 || (UserTransaction.TotalTransactions != 1 && UserTransaction.MinAmount > transaction.Amount) {
		UserTransaction.MinAmount = transaction.Amount
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("transaction inserted")

}

// GetStatics - provides statics values of transactions
func GetStatics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")

	// get current user location from request header
	currentLocation := r.Header.Get("location")

	// check location of user is valid
	if UserTransaction.City.City != "" && currentLocation != UserTransaction.City.City {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if UserTransaction.TotalTransactions == 0 {
		http.Error(w, "transactions not found", http.StatusNoContent)
		return
	}

	// New struct object of Statics to store response data
	var statics = Statics{
		Sum:   UserTransaction.TotalAmount,
		Avg:   UserTransaction.TotalAmount / float64(UserTransaction.TotalTransactions),
		Max:   UserTransaction.MaxAmount,
		Min:   UserTransaction.MinAmount,
		Count: UserTransaction.TotalTransactions,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(statics)
}

// DeleteTransaction - delete all transaction from the transactions list.
func DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")

	UserTransaction.Transactions = []Transaction{}
	UserTransaction.MaxAmount = 0
	UserTransaction.MinAmount = 0
	UserTransaction.TotalAmount = 0
	UserTransaction.TotalTransactions = 0

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode("transaction deleted")
}

// SetUserCity - update the city of end user
func SetUserCity(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")

	var city City
	if err := json.NewDecoder(r.Body).Decode(&city); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	// remove white spaces in both end
	city.City = strings.TrimSpace(city.City)

	if city.City == "" {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	UserTransaction.City = city

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("location updated")
}

// ResetUserCity - reset the city of end user. Value changes to an empty string.
func ResetUserCity(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")

	UserTransaction.City.City = ""

	w.WriteHeader(http.StatusResetContent)
	json.NewEncoder(w).Encode("location reset")
}
