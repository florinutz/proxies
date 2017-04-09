package proxies

import "net"

const (
	proxyListURL = "http://txt.proxyspy.net/proxy.txt"
	myIpUrl      = "https://api.ipify.org?format=json"
)

type Proxy struct {
	IP           net.IPAddr
	Port         int
	Country      string
	Anonymity    string
	WithTLS      bool
	GooglePassed bool
}

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
