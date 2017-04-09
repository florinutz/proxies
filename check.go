package proxies

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"time"
)

var Transport = &http.Transport{
	DialContext: (&net.Dialer{
		Timeout:   5 * time.Second,
		DualStack: true,
	}).DialContext,
	DisableKeepAlives:     true,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
	TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
}

var Client = &http.Client{
	Transport: Transport,
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		// As a special case, if CheckRedirect returns ErrUseLastResponse,
		// then the most recent response is returned with its body unclosed,
		// along with a nil error.
		return http.ErrUseLastResponse
	},
}

// Performs a GET request to a targetUrl through the proxy and returns the response
func Check(proxy string, targetUrl string) (*http.Response, error) {
	proxyURL, err := url.Parse(proxy)
	if err != nil {
		return nil, err
	}
	Transport.Proxy = http.ProxyURL(proxyURL)

	res, err := Client.Get(targetUrl)
	if err != nil {
		return nil, err
	}

	return res, nil
}
