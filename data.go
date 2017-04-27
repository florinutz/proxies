package proxies

import (
	"fmt"
	"io"
	"net"
	"net/http"
)

//https://github.com/Kladdkaka/ProxyGrabber/blob/master/proxy_grabber.py

const (
	proxyListURL = "http://txt.proxyspy.net/proxy.txt"
	myIpUrl      = "https://api.ipify.org?format=json"
)

// Extracts a list of proxies from io.Reader
type ListReader interface {
	Read(reader io.Reader) *List
}

type ListReaderFunc func(reader io.Reader) *List

func (f ListReaderFunc) Read(r io.Reader) *List {
	return f(r)
}

// Matches http.Client, allows changing the request
// for timeouts, mainly
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

type List []Proxy

type Proxy struct {
	IP           net.IPAddr
	Port         int
	Country      string
	Anonymity    string
	TLS          string
	GooglePassed bool
}

func (p *Proxy) String() string {
	return fmt.Sprintf("%s:%d", p.IP, p.Port)
}

func (l *List) Read(reader io.Reader, f ListReader) {
	*l = *f.Read(reader)
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
