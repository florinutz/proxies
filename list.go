package proxies

import (
	"bufio"
	"io"
	"net"
	"strconv"
	"strings"
)

type List struct {
	Parser  ListParser
	Proxies []Proxy
}

type ListParser struct {
	reader io.Reader
}

func NewListParser(r io.Reader) *ListParser {
	return &ListParser{r}
}

func (parser *ListParser) GetProxies() (proxies []Proxy) {
	proxies = []Proxy{}
	scanner := bufio.NewScanner(parser.reader)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 || '0' > line[0] || line[0] > '9' {
			continue
		}
		address := strings.Fields(line)[0]
		addressSplit := strings.Split(address, ":")

		port, err := strconv.Atoi(addressSplit[1])
		if err != nil {
			// or a default value?
			continue
		}

		proxy := Proxy{
			IP: net.IPAddr{
				IP: []byte(addressSplit[0]),
			},
			Port: port,
		}
		proxies = append(proxies, proxy)
	}
	return
}
