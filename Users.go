package coinbaseDriver

import (
	"net/http"
	"time"
)

type Userinformation struct {
	ID              string      `json:"id"`
	Name            string      `json:"name"`
	Username        interface{} `json:"username"`
	ProfileLocation interface{} `json:"profile_location"`
	ProfileBio      interface{} `json:"profile_bio"`
	ProfileURL      interface{} `json:"profile_url"`
	AvatarURL       string      `json:"avatar_url"`
	Resource        string      `json:"resource"`
	ResourcePath    string      `json:"resource_path"`
	Email           string      `json:"email"`
	TimeZone        string      `json:"time_zone"`
	NativeCurrency  string      `json:"native_currency"`
	BitcoinUnit     string      `json:"bitcoin_unit"`
	State           interface{} `json:"state"`
	Country         struct {
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"country"`
	CreatedAt time.Time `json:"created_at"`
}

// UserInformation returns informations about the user
func (client *CoinbaseClient) UserInformation() (interface{}, error) {
	var typ Userinformation
	err := client.query(&typ, http.MethodGet, "user", nil)
	return typ, err
}
