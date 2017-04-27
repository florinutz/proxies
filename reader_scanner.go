package proxies

import (
	"bufio"
	"io"
	"net"
	"regexp"
	"strconv"
)

type ProxyExtractor interface {
	ExtractProxy(string) *Proxy
}

// Implements the ProxyExtractor interface
type ProxyExtractorFunc func(string) *Proxy

func (f ProxyExtractorFunc) ExtractProxy(s string) *Proxy {
	return f(s)
}

type ScannerReader struct {
	ProxyExtractor ProxyExtractor
	SplitFunc      bufio.SplitFunc
}

// Implements the ListReader interface
func (r *ScannerReader) Read(reader io.Reader) *List {
	list := List{}
	scanner := bufio.NewScanner(reader)
	if r.SplitFunc != nil {
		scanner.Split(r.SplitFunc)
	}
	for scanner.Scan() {
		txt := scanner.Text()
		if proxy := r.ProxyExtractor.ExtractProxy(txt); proxy != nil {
			list = append(list, *proxy)
		}
	}
	return &list
}

// Returns *Proxy with populated values or nil on failure
func readSingleProxy(s string) *Proxy {
	var lineSplitterRegex = `(?P<hostport>\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}:\d{1,5})\s(?P<country>\w{1,3})-(?P<anonimity>[NAH]!?)(?:-(?P<tls>S!?))?\s(?P<google>[+-])`
	values := getLineNamedValues(regexp.MustCompile(lineSplitterRegex), s)

	ip, portStr, err := net.SplitHostPort(values["hostport"])
	if err != nil {
		return nil
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil
	}
	google := false
	if values["google"] == "+" {
		google = true
	}

	return &Proxy{
		IP:           net.IPAddr{IP: net.ParseIP(ip)},
		Port:         port,
		Country:      values["country"],
		Anonymity:    values["anonimity"],
		TLS:          values["tls"],
		GooglePassed: google,
	}
}

func getLineNamedValues(re *regexp.Regexp, line string) (values map[string]string) {
	values = map[string]string{}
	if match := re.FindStringSubmatch(line); match != nil {
		n1 := re.SubexpNames()
		for i, n := range match {
			values[n1[i]] = n
		}
	}
	return
}
