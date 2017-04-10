package proxies

import (
	"bufio"
	"io"
	"net"
	"regexp"
	"strconv"
)

type List struct {
	Proxies []*Proxy
}

func (l *List) Read(reader io.Reader) (proxies []*Proxy) {
	proxies = []*Proxy{}
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		txt := scanner.Text()
		if len(txt) == 0 || '0' > txt[0] || txt[0] > '9' {
			continue
		}
		if proxy := l.ReadSingleProxy(txt); proxy != nil {
			proxies = append(proxies, proxy)
		}
	}
	return
}

// Returns *Proxy with populated values or nil on failure
func (l *List) ReadSingleProxy(s string) *Proxy {
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
