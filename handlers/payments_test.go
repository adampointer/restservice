package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/adampointer/restservice/data"
	"github.com/gorilla/mux"
)

func getTestDB(t *testing.T) *data.Client {
	dir, err := ioutil.TempDir("", "handler_tests")
	if err != nil {
		t.Fatalf("unable to create test db: %s", err)
	}
	dbClient, err := data.NewClient(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("unable to create test db: %s", err)
	}
	return dbClient
}

func cleanUp(db *data.Client) {
	db.Close()
	// Ignore errors here
	os.Remove(db.Path())
}

func TestGetAllPaymentsEmpty(t *testing.T) {
	db := getTestDB(t)
	defer cleanUp(db)
	h := NewPayments(db)

	req, err := http.NewRequest("GET", "/payments", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetAll)
	handler.ServeHTTP(rr, req)

	// Assert that a GET to /payments when there are no payments returns a 200 and []
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got '%v' want '%v'", status, http.StatusOK)
	}
	expected := "[]"
	actual := strings.TrimRight(rr.Body.String(), "\r\n")
	if actual != expected {
		t.Errorf("handler returned unexpected body: got '%s' want '%s'", actual, expected)
	}
}

func TestCreateThenGetAllPayments(t *testing.T) {
	id := "4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43"

	db := getTestDB(t)
	defer cleanUp(db)
	h := NewPayments(db)

	// First, assert a unique resource is created
	req, err := http.NewRequest("PUT", "/payments/"+id, strings.NewReader(exampleJSON))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/payments/{id}", h.Create)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got '%v' want '%v'", status, http.StatusCreated)
	}

	// Second, assert that calling GetAll returns the new resource in an array of one element
	req, err = http.NewRequest("GET", "/payments", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetAll)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got '%v' want '%v'", status, http.StatusOK)
	}
	var payments []*data.Payment
	if err := json.NewDecoder(rr.Body).Decode(&payments); err != nil {
		t.Fatal("unable to decode response into JSON")
	}
	if len(payments) != 1 {
		t.Fatalf("handler returned %d payments, we wanted 1", len(payments))
	}
	// Choose a couple of fields to match here - we could do a deep compare but it's probably overkill
	if payments[0].ID != id {
		t.Fatalf("handler returned payment with id '%s', we expected '%s'", payments[0].ID, id)
	}
	if payments[0].Attributes.BeneficiaryParty.AccountName != "W Owens" {
		t.Fatalf("handler returned payment with beneficiary account name '%s', we expected 'W Owens'",
			payments[0].Attributes.BeneficiaryParty.AccountName)
	}
}

func TestCreateDuplicateID(t *testing.T) {}
func TestCreateNoID(t *testing.T)        {}
func TestCreateEmptyBody(t *testing.T)   {}

var exampleJSON = `{
	"type": "Payment",
	"id": "",
	"version": 0,
	"organisation_id": "743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb",
	"attributes": {
		"amount": "100.21",
		"beneficiary_party": {
			"account_name": "W Owens",
			"account_number": "31926819",
			"account_number_code": "BBAN",
			"account_type": 0,
			"address": "1 The Beneficiary Localtown SE2",
			"bank_id": "403000",
			"bank_id_code": "GBDSC",
			"name": "Wilfred Jeremiah Owens"
		},
		"charges_information": {
			"bearer_code": "SHAR",
			"sender_charges": [{
				"amount": "5.00",
				"currency": "GBP"
			}, {
				"amount": "10.00",
				"currency": "USD"
			}],
			"receiver_charges_amount": "1.00",
			"receiver_charges_currency": "USD"
		},
		"currency": "GBP",
		"debtor_party": {
			"account_name": "EJ Brown Black",
			"account_number": "GB29XABC10161234567801",
			"account_number_code": "IBAN",
			"address": "10 Debtor Crescent Sourcetown NE1",
			"bank_id": "203301",
			"bank_id_code": "GBDSC",
			"name": "Emelia Jane Brown"
		},
		"end_to_end_reference": "Wil piano Jan",
		"fx": {
			"contract_reference": "FX123",
			"exchange_rate": "2.00000",
			"original_amount": "200.42",
			"original_currency": "USD"
		},
		"numeric_reference": "1002001",
		"payment_id": "123456789012345678",
		"payment_purpose": "Paying for goods/services",
		"payment_scheme": "FPS",
		"payment_type": "Credit",
		"processing_date": "2017-01-18",
		"reference": "Payment for Em's piano lessons",
		"scheme_payment_sub_type": "InternetBanking",
		"scheme_payment_type": "ImmediatePayment",
		"sponsor_party": {
			"account_number": "56781234",
			"bank_id": "123123",
			"bank_id_code": "GBDSC"
		}
	}
}`
