package validator

import (
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"time"
)

const (
	myIPURL = "https://api.ipify.org?format=json"
)

// Check verifies a proxy
func Check(proxy string) (bool, error) {
	proxyURL, err := url.Parse(proxy)
	if err != nil {
		return false, err
	}

	tp := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			DualStack: true,
		}).DialContext,
		DisableKeepAlives:     true,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	client := &http.Client{Transport: tp}

	res, err := client.Get(myIPURL)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	var jr struct {
		IP string `json:"ip"`
	}

	err = json.NewDecoder(res.Body).Decode(&jr)
	if err != nil {
		return false, err
	}

	if len(jr.IP) == 0 {
		return false, nil
	}

	return false, nil
}
