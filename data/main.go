package data

import (
	"fmt"

	"github.com/asdine/storm"
)

// Client abstracts our database
type Client struct {
	dbPath string
	db     *storm.DB
}

// NewClient returns a new client with database at path
func NewClient(path string) (*Client, error) {
	db, err := storm.Open(path)
	if err != nil {
		return nil, err
	}
	c := &Client{
		dbPath: path,
		db:     db,
	}
	return c, nil
}

// Close the database
func (c *Client) Close() {
	c.db.Close()
}

// Path returns the db path
func (c *Client) Path() string {
	return c.dbPath
}

// FetchPayment gets a single Payment by ID
func (c *Client) FetchPayment(id string) (*Payment, error) {
	var pmt Payment
	if err := c.db.One("ID", id, &pmt); err != nil {
		return nil, err
	}
	return &pmt, nil
}

// FetchAllPayments gets a single Payment by ID
func (c *Client) FetchAllPayments() ([]*Payment, error) {
	var pmts []*Payment
	if err := c.db.All(&pmts); err != nil {
		return nil, err
	}
	return pmts, nil
}

// CreatePayment saves a new Payment in the database
func (c *Client) CreatePayment(pmt *Payment) error {
	_, err := c.FetchPayment(pmt.ID)
	if err == nil || err.Error() != "not found" {
		return fmt.Errorf("resource exists")
	}
	return c.db.Save(pmt)
}

// UpdatePayment updates an existing Payment in the database
func (c *Client) UpdatePayment(pmt *Payment) error {
	return c.db.Update(pmt)
}

// DeletePayment deletes an existing Payment from the database
func (c *Client) DeletePayment(id string) error {
	pmt, err := c.FetchPayment(id)
	if err != nil {
		return err
	}
	return c.db.DeleteStruct(pmt)
}
