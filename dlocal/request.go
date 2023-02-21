package dlocal

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"os"
	"time"
)

func PostRequest(body []byte, endpoint string) (*http.Request, error) {
	host := os.Getenv("DLOCAL_URL")
	x_login := os.Getenv("DLOCAL_X_LOGIN")
	x_trans_key := os.Getenv("DLOCAL_X_TRANS_KEY")
	x_date := time.Now().Format(time.RFC3339)

	req, err := http.NewRequest(http.MethodPost, host+endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Date", x_date)
	req.Header.Set("X-Login", x_login)
	req.Header.Set("X-Trans-Key", x_trans_key)

	// Authorization Header
	secret := os.Getenv("DLOCAL_SECRET")

	// Create a new HMAC by defining the hash type and the key (as byte array)
	h := hmac.New(sha256.New, []byte(secret))

	data := []byte(x_login + x_date)
	// Write Data to it
	h.Write(append(data, body...))

	// Get result and encode as hexadecimal string
	sha := hex.EncodeToString(h.Sum(nil))
	// Set Header
	req.Header.Set("Authorization", "V2-HMAC-SHA256, Signature: "+sha)

	return req, nil
}
