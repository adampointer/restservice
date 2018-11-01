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

func TestGetPaymentNotFound(t *testing.T) {
	db := getTestDB(t)
	defer cleanUp(db)
	h := NewPayments(db)

	req, err := http.NewRequest("GET", "/payments/foobar", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/payments/{id}", h.GetOne)
	router.ServeHTTP(rr, req)

	// Assert that a GET to /payments/{ID} with an invalid id returns 404
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got '%v' want '%v'", status, http.StatusNotFound)
	}
}

func TestCreatePaymentThenGetAllPayments(t *testing.T) {
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

func TestCreatePaymentThenGetPayment(t *testing.T) {
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

	// Second, assert that calling Get returns the new resource
	req, err = http.NewRequest("GET", "/payments/"+id, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	router = mux.NewRouter()
	router.HandleFunc("/payments/{id}", h.GetOne)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got '%v' want '%v'", status, http.StatusOK)
	}
	var payment data.Payment
	if err := json.NewDecoder(rr.Body).Decode(&payment); err != nil {
		t.Fatal("unable to decode response into JSON")
	}
	if payment.ID != id {
		t.Fatalf("handler returned payment with id '%s', we expected '%s'", payment.ID, id)
	}
	if payment.Attributes.BeneficiaryParty.AccountName != "W Owens" {
		t.Fatalf("handler returned payment with beneficiary account name '%s', we expected 'W Owens'",
			payment.Attributes.BeneficiaryParty.AccountName)
	}
}

func TestCreatePaymentDuplicateID(t *testing.T) {
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

	// Next repeat and check it was a 400
	req, err = http.NewRequest("PUT", "/payments/"+id, strings.NewReader(exampleJSON))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	router = mux.NewRouter()
	router.HandleFunc("/payments/{id}", h.Create)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got '%v' want '%v'", status, http.StatusBadRequest)
	}
}

func TestCreatePaymentNoID(t *testing.T) {
	db := getTestDB(t)
	defer cleanUp(db)
	h := NewPayments(db)

	// See what happens if don't include an ID - actually this should not be possible in the real world
	req, err := http.NewRequest("PUT", "/payments/+", strings.NewReader(exampleJSON))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.Create)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got '%v' want '%v'", status, http.StatusBadRequest)
	}
}

func TestCreatePaymentEmptyBody(t *testing.T) {
	db := getTestDB(t)
	defer cleanUp(db)
	h := NewPayments(db)

	// See what happens if we PUT without a body
	req, err := http.NewRequest("PUT", "/payments/+", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.Create)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got '%v' want '%v'", status, http.StatusBadRequest)
	}

}

func TestUpdatePaymentHappyPath(t *testing.T) {
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

	// Next update and check we get 200
	req, err = http.NewRequest("POST", "/payments/"+id, strings.NewReader(exampleJSON2))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	router = mux.NewRouter()
	router.HandleFunc("/payments/{id}", h.Update)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got '%v' want '%v'", status, http.StatusOK)
	}

	// Check our change has been made
	req, err = http.NewRequest("GET", "/payments/"+id, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	router = mux.NewRouter()
	router.HandleFunc("/payments/{id}", h.GetOne)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got '%v' want '%v'", status, http.StatusOK)
	}
	var payment data.Payment
	if err := json.NewDecoder(rr.Body).Decode(&payment); err != nil {
		t.Fatal("unable to decode response into JSON")
	}
	if payment.ID != id {
		t.Fatalf("handler returned payment with id '%s', we expected '%s'", payment.ID, id)
	}
	if payment.Attributes.BeneficiaryParty.AccountName != "Foo Bar" {
		t.Fatalf("handler returned payment with beneficiary account name '%s', we expected 'Foo Bar'",
			payment.Attributes.BeneficiaryParty.AccountName)
	}
}

func TestUpdatePaymentBadRequest(t *testing.T) {
	db := getTestDB(t)
	defer cleanUp(db)
	h := NewPayments(db)

	req, err := http.NewRequest("POST", "/payments/foobar", strings.NewReader(exampleJSON2))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/payments/{id}", h.Update)
	router.ServeHTTP(rr, req)

	// Assert that a POST to /payments/{ID} with an invalid id returns 404
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got '%v' want '%v'", status, http.StatusNotFound)
	}
}

func TestDeletePaymentHappyPath(t *testing.T) {
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

	// Next delete and check we get 200
	req, err = http.NewRequest("DELETE", "/payments/"+id, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	router = mux.NewRouter()
	router.HandleFunc("/payments/{id}", h.Delete)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got '%v' want '%v'", status, http.StatusOK)
	}

	// Check our payment is gone
	req, err = http.NewRequest("GET", "/payments/"+id, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	router = mux.NewRouter()
	router.HandleFunc("/payments/{id}", h.GetOne)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got '%v' want '%v'", status, http.StatusNotFound)
	}
}

func TestDeletePaymentNotFound(t *testing.T) {
	db := getTestDB(t)
	defer cleanUp(db)
	h := NewPayments(db)

	req, err := http.NewRequest("DELETE", "/payments/foobar", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/payments/{id}", h.Delete)
	router.ServeHTTP(rr, req)

	// Assert that a DELETE to /payments/{ID} with an invalid id returns 404
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got '%v' want '%v'", status, http.StatusNotFound)
	}
}
