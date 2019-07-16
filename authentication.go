package upvest

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"io/ioutil"
	"time"
)

// Headers represent the HTTP headers sent to Upvest API
type Headers map[string]string

// AuthProvider interface for authentication mechanisms supported by Upvest API
type AuthProvider interface {
	// GetHeaders returns authorization headers (or other info) to be attached to requests.
	GetHeaders(method, path string, body interface{}) (Headers, error)
}

// KeyAuth (The API Key Authentication) is used to authenticate requests as a tenant.
type KeyAuth struct {
	apiKey        string
	apiSecret     string
	apiPassphrase string
}

// GetHeaders returns authorization headers for requests as a tenant.
func (auth KeyAuth) GetHeaders(method, path string, body interface{}) (Headers, error) {
	var headers Headers
	timestamp := string(time.Now().UnixNano() / 1000000)
	// Compose the message as a concatenation of all info we are sending along with the request
	message := timestamp + method + apiVersion + path

	if body != nil {
		buf, err := jsonEncode(body)
		if err != nil {
			return nil, err
		}
		body1, _ := ioutil.ReadAll(buf)
		message = message + string(body1)
	}

	// Generate signature, in order to prevent manipulation of payload in flight
	h := hmac.New(sha512.New, []byte(message))
	signature := hex.EncodeToString(h.Sum(nil))

	// Generate message headers
	headers = Headers{
		"Content-Type":         "application/json",
		"X-UP-API-Key":         auth.apiKey,
		"X-UP-API-Signature":   signature,
		"X-UP-API-Timestamp":   timestamp,
		"X-UP-API-Passphrase":  auth.apiPassphrase,
		"X-UP-API-Signed-Path": apiVersion + path,
	}

	return headers, nil
}