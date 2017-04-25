package proxies

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"regexp"
	"strconv"
)

//https://github.com/Kladdkaka/ProxyGrabber/blob/master/proxy_grabber.py

const (
	proxyListURL = "http://txt.proxyspy.net/proxy.txt"
	myIpUrl      = "https://api.ipify.org?format=json"
)

type Proxy struct {
	IP           net.IPAddr
	Port         int
	Country      string
	Anonymity    string
	TLS          string
	GooglePassed bool
}

type List []Proxy

func (p *Proxy) String() string {
	return fmt.Sprintf("%s:%d", p.IP, p.Port)
}

func (l *List) Read(reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		txt := scanner.Text()
		if len(txt) == 0 || '0' > txt[0] || txt[0] > '9' {
			continue
		}
		if proxy := l.ReadSingleProxy(txt); proxy != nil {
			*l = append(*l, *proxy)
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

// here just for the waitingGroup:
//
//func main() {
//	res, err := http.Get(proxyListURL)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Checking %s\nHitting %s\n\n", color.GreenString(proxyListURL), color.GreenString(myIpUrl))
//	if res.Status != "200 OK" {
//		log.Fatalf("Wrong status in response: %s", res.Status)
//	}
//	defer res.Body.Close()
//
//	var wg sync.WaitGroup
//
//	scn := bufio.NewScanner(res.Body)
//	for scn.Scan() {
//		line := scn.Text()
//		if len(line) == 0 || '0' > line[0] || line[0] > '9' {
//			continue
//		}
//
//		parts := strings.Fields(line)
//		proxy := "http://" + parts[0]
//
//		wg.Add(1)
//		go func(proxy string) {
//			defer wg.Done()
//			_, err := checker.Check(proxy, myIpUrl)
//			if err != nil {
//				log.Printf("proxy %s doesn't work: %s", proxy, color.RedString(err.Error()))
//			} else {
//				log.Printf("proxy %s %s", proxy, color.GreenString("works"))
//			}
//		}(proxy)
//	}
//	if err := scn.Err(); err != nil {
//		log.Fatal(err)
//	}
//
//	wg.Wait()
//}
