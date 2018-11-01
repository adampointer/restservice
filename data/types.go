package data

import "encoding/json"

// Resource contains the base attributes
type Resource struct {
	Type           string `json:"type"`
	ID             string `json:"id",storm:"id"`
	Version        int    `json:"version"`
	OrganisationID string `json:"organisation_id"`
}

// Payment represents a payment resource
type Payment struct {
	Resource   `storm:"inline"`
	Attributes *PaymentAttributes `json:"attributes"`
}

// PaymentAttributes are the actual details of the payment
type PaymentAttributes struct {
	Amount               json.Number     `json:"amount"`
	BeneficiaryParty     *PaymentParty   `json:"beneficiary_party"`
	ChargesInformation   *PaymentCharges `json:"charges_information"`
	Currency             string          `json:"currency"`
	DebtorParty          *PaymentParty   `json:"debtor_party"`
	EtoEReference        string          `json:"end_to_end_reference"`
	FX                   *PaymentFXData  `json:"fx"`
	NumericReference     json.Number     `json:"numeric_reference",storm:"unique"`
	PaymentID            json.Number     `json:"payment_id",storm:"unique"`
	PaymentPurpose       string          `json:"payment_purpose"`
	PaymentScheme        string          `json:"payment_scheme"`
	PaymentType          string          `json:"payment_type"`
	ProcessingDate       string          `json:"processing_date"`
	Reference            string          `json:"reference"`
	SchemePaymentSubType string          `json:"scheme_payment_sub_type"`
	SchemePaymentType    string          `json:"scheme_payment_type"`
	SponsorParty         *PaymentParty   `json:"sponsor_party"`
}

// PaymentParty represents a party involved in the transaction
type PaymentParty struct {
	AccountName       string `json:"account_name"`
	AccountNumber     string `json:"account_number"`
	AccountNumberCode string `json:"account_number_code"`
	AccountType       int    `json:"account_type"`
	Address           string `json:"address"`
	BankID            string `json:"bank_id"`
	BankIDCode        string `json:"bank_id_code"`
	Name              string `json:"name"`
}

// PaymentCharges is the charges associated with the payment
type PaymentCharges struct {
	BearerCode              string                 `json:"bearer_code"`
	SenderCharges           []*PaymentSenderCharge `json:"sender_charges"`
	ReceiverChargesAmount   json.Number            `json:"receiver_charges_amount"`
	ReceiverChargesCurrency string                 `json:"receiver_charges_currency"`
}

// PaymentSenderCharge is a charge on the payment sender
type PaymentSenderCharge struct {
	Amount   json.Number `json:"amount"`
	Currency string      `json:"currency"`
}

// PaymentFXData is the forex details of the payment
type PaymentFXData struct {
	ContractReference string      `json:"contract_reference"`
	ExchangeRate      json.Number `json:"exchange_rate"`
	OriginalAmount    json.Number `json:"original_amount"`
	OriginalCurrency  string      `json:"original_currency"`
}
