package coinbaseDriver

import (
	"net/http"
)

// Price represent a sell price
type Price struct {
	Base     string `json:"base"`
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
}

// Prices returns prices for an pair (for example: BTC-USD)
func (client *CoinbaseClient) Prices(pair string) (Price, error) {
	var price Price
	err := client.query(&price, http.MethodGet, "prices/"+pair+"/buy", nil)
	return price, err
}
