package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/adampointer/restservice/data"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// Payments handlers for payment resorces
type Payments struct {
	db *data.Client
}

// ErrNotFound - resource not found 404
var ErrNotFound = errors.New("not found")

// NewPayments returns new handler with database client
func NewPayments(db *data.Client) *Payments {
	return &Payments{db: db}
}

// GetAll lists all payment resources
func (p *Payments) GetAll(w http.ResponseWriter, r *http.Request) {
	pmts, err := p.db.FetchAllPayments()
	if err != nil {
		log.Errorf("Error getting all payments: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(pmts)
}

// GetOne shows a single payment resource
func (p *Payments) GetOne(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	pmt, err := p.db.FetchPayment(params["id"])
	if err != nil {
		if err.Error() == ErrNotFound.Error() {
			w.WriteHeader(http.StatusNotFound)
		} else {
			log.Errorf("Error getting payments: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	json.NewEncoder(w).Encode(pmt)
}

// Create a new payment resource
func (p *Payments) Create(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var payment data.Payment
	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		log.Errorf("Error decoding create payment request: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	payment.ID = params["id"]
	if len(payment.ID) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := p.db.CreatePayment(&payment); err != nil {
		if err.Error() == "resource exists" {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			log.Errorf("Error saving new payment: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// Update an existing payment resource
func (p *Payments) Update(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var payment data.Payment
	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		log.Errorf("Error decoding update payment request: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	payment.ID = params["id"]
	if err := p.db.UpdatePayment(&payment); err != nil {
		if err.Error() == ErrNotFound.Error() {
			w.WriteHeader(http.StatusNotFound)
		} else {
			log.Errorf("Error saving new payment: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Delete a payment resource
func (p *Payments) Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if err := p.db.DeletePayment(params["id"]); err != nil {
		if err.Error() == ErrNotFound.Error() {
			w.WriteHeader(http.StatusNotFound)
		} else {
			log.Errorf("Error deleting payment: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}
