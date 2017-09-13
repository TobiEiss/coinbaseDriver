package coinbaseDriver

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// CoinbaseClient represent your client
type CoinbaseClient struct {
	context.Context
	Host         string
	AccessKey    string
	AccessSecret string
}

// CoinbaseResponseWrapper wrapps the data
type CoinbaseResponseWrapper struct {
	Typ interface{} `json:"data"`
}

// HTTPDo function runs the HTTP request and processes its response in a new goroutine.
func HTTPDo(ctx context.Context, request *http.Request, processResponse func(*http.Response, error) error) error {
	// Run the HTTP request in a goroutine and pass the response to processResponse.
	transport := &http.Transport{}
	client := &http.Client{Transport: transport}
	errorChannel := make(chan error, 1)

	// do request
	go func() { errorChannel <- processResponse(client.Do(request)) }()
	select {
	case <-ctx.Done():
		transport.CancelRequest(request)
		<-errorChannel // wait for processResponse function
		return ctx.Err()
	case err := <-errorChannel:
		return err
	}
}

// query the api
func (coinbaseClient *CoinbaseClient) query(typ interface{}, reqMeth string, route string, values url.Values) error {
	var coinbaseResponse CoinbaseResponseWrapper
	coinbaseResponse.Typ = typ

	// create httpURL
	httpURL := coinbaseClient.Host + route

	// create http-Context
	httpContext, cancelFunc := context.WithTimeout(coinbaseClient, 15*time.Second)
	defer cancelFunc()

	// build request
	var reader io.Reader
	if values != nil {
		reader = strings.NewReader(values.Encode())
	}
	request, err := http.NewRequest(reqMeth, httpURL, reader)
	if err != nil {
		return err
	}

	// add headers
	path := html.EscapeString(request.URL.Path)
	timestamp := fmt.Sprintf("%d", time.Now().UTC().Unix())

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("CB-VERSION", "2017-09-13")
	request.Header.Set("CB-ACCESS-KEY", coinbaseClient.AccessKey)
	request.Header.Set("CB-ACCESS-TIMESTAMP", timestamp)
	request.Header.Set("CB-ACCESS-SIGN", coinbaseClient.sign(timestamp, request.Method, path, ""))

	// fire up request and unmarshal serverTime
	err = HTTPDo(httpContext, request, func(response *http.Response, err error) error {
		if err != nil {
			return err
		}
		defer response.Body.Close()

		decoder := json.NewDecoder(response.Body)
		if err := decoder.Decode(&coinbaseResponse); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (coinbaseClient *CoinbaseClient) sign(timestamp, method, path, body string) string {
	h := hmac.New(sha256.New, []byte(coinbaseClient.AccessSecret))
	h.Write([]byte(timestamp + method + path + body))
	return hex.EncodeToString(h.Sum(nil))
}

func createHmac(timestamp string, method string, requestPath string, body string, secret string) string {
	message := timestamp + method + requestPath + body
	hmac := hmac.New(sha1.New, []byte(secret))
	hmac.Write([]byte(message))
	return hex.EncodeToString(hmac.Sum(nil))
}
