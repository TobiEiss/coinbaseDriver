package coinbaseDriver

import (
	"net/http"
	"net/url"
	"time"
)

// Accounts returns all accounts
// TODO: build a account-model
func (client *CoinbaseClient) Accounts() (interface{}, error) {
	var accounts interface{}
	err := client.query(&accounts, http.MethodGet, "accounts", nil)
	return accounts, err
}

type Receipt struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Status string `json:"status"`
	Amount struct {
		Amount   string `json:"amount"`
		Currency string `json:"currency"`
	} `json:"amount"`
	NativeAmount struct {
		Amount   string `json:"amount"`
		Currency string `json:"currency"`
	} `json:"native_amount"`
	Description     interface{} `json:"description"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
	Resource        string      `json:"resource"`
	ResourcePath    string      `json:"resource_path"`
	InstantExchange bool        `json:"instant_exchange"`
	Network         struct {
		Status         string `json:"status"`
		TransactionFee struct {
			Amount   string `json:"amount"`
			Currency string `json:"currency"`
		} `json:"transaction_fee"`
		TransactionAmount struct {
			Amount   string `json:"amount"`
			Currency string `json:"currency"`
		} `json:"transaction_amount"`
		Confirmations int `json:"confirmations"`
	} `json:"network"`
	To struct {
		Resource string `json:"resource"`
		Address  string `json:"address"`
	} `json:"to"`
	Details struct {
		Title    string `json:"title"`
		Subtitle string `json:"subtitle"`
	} `json:"details"`
}

// SendMoney sends money to a wallet-address for any currency
func (client *CoinbaseClient) SendMoney(to string, amount string, currency string) (interface{}, error) {
	var receipt Receipt
	values := url.Values{}
	values.Add("type", "send")
	values.Add("to", to)
	values.Add("amount", amount)
	values.Add("currency", currency)
	err := client.query(&receipt, http.MethodPost, "accounts/8700de83-86df-5ba7-b321-f48fdc1e92a9/transactions", values)
	return receipt, err
}
