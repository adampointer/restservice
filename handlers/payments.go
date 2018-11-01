package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/adampointer/restservice/data"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// Payments handlers for payment resorces
type Payments struct {
	db *data.Client
}

// NewPayments returns new handler with database client
func NewPayments(db *data.Client) *Payments {
	return &Payments{db: db}
}

// GetAll lists all payment resources
func (p *Payments) GetAll(w http.ResponseWriter, r *http.Request) {}

// GetOne shows a single payment resource
func (p *Payments) GetOne(w http.ResponseWriter, r *http.Request) {}

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
	if err := p.db.CreatePayment(&payment); err != nil {
		log.Errorf("Error saving new payment: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// Update an existing payment resource
func (p *Payments) Update(w http.ResponseWriter, r *http.Request) {}

// Delete a payment resource
func (p *Payments) Delete(w http.ResponseWriter, r *http.Request) {}
