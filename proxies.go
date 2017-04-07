package main

import (
	"bufio"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/florinutz/proxies/validator"
)

const (
	proxyListURL = "http://txt.proxyspy.net/proxy.txt"
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
			ok, err := validator.Check(proxy)
			if err != nil {
				return
			}
			if ok {
				log.Printf("proxy %s works", proxy)
			} else {
				log.Printf("proxy %s doesn't work", proxy)
			}
		}(proxy)

	}
	if err := scn.Err(); err != nil {
		log.Fatal(err)
	}

	wg.Wait()
}
