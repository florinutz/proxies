package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	proxyListURL = "http://txt.proxyspy.net/proxy.txt"
	myIPURL      = "https://api.ipify.org?format=json"
)

func main() {
	res, err := http.Get(proxyListURL)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	var wg sync.WaitGroup

	scn := bufio.NewScanner(res.Body)
	for scn.Scan() {
		line := scn.Text()
		if len(line) == 0 || '0' > line[0] || line[0] > '9' {
			continue
		}

		parts := strings.Fields(line)
		proxy := "http://" + parts[0]

		wg.Add(1)
		go func(proxy string) {
			defer wg.Done()
			err := checkProxy(proxy)
			if err != nil {
				return
			}
			log.Printf("proxy %s works", proxy)
		}(proxy)

	}
	if err := scn.Err(); err != nil {
		log.Fatal(err)
	}

	wg.Wait()
}

func checkProxy(proxy string) error {
	proxyURL, err := url.Parse(proxy)
	if err != nil {
		return err
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

	client := &http.Client{
		Transport: tp,
	}

	res, err := client.Get(myIPURL)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var jr struct {
		IP string `json:"ip"`
	}

	err = json.NewDecoder(res.Body).Decode(&jr)
	if err != nil {
		return err
	}

	if len(jr.IP) == 0 {
		return errors.New("proxy doesn't work")
	}

	return nil
}
