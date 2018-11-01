package data

import "github.com/asdine/storm"

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
	return c.db.Save(pmt)
}

// UpdatePayment updates an existing Payment in the database
func (c *Client) UpdatePayment(pmt *Payment) error {
	return c.db.Update(pmt)
}

// DeletePayment deletes an existing Payment from the database
func (c *Client) DeletePayment(pmt *Payment) error {
	return c.db.DeleteStruct(pmt)
}
