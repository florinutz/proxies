package proxies

import (
	"bufio"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestReadFile(t *testing.T) {
	t.Parallel()

	list := List{}

	f, err := os.Open("test_data/list_sample.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	list.Read(r, &ScannerReader{
		ProxyExtractor: ProxyExtractorFunc(readSingleProxy),
	})

	// there are 300 possibly valid proxiesRead in list_sample.txt
	if len(list) != 300 {
		t.Fatal("Not 300 proxiesRead")
	}
}

func TestReadLine(t *testing.T) {
	t.Parallel()

	proxyExtractor := ProxyExtractorFunc(readSingleProxy)

	if proxy := proxyExtractor("54.153.73.116:3128 US-N-S! -"); proxy != nil {
		if proxy.IP.String() != "54.153.73.116" {
			t.Error("Wrong ip")
		}
		if proxy.Port != 3128 {
			t.Error("Wrong port")
		}
		if proxy.Country != "US" {
			t.Error("Wrong country")
		}
		if proxy.Anonymity != "N" {
			t.Error("Wrong anonimity type")
		}
		if proxy.TLS != "S!" {
			t.Error("Wrong TLS capabilities")
		}
		if proxy.GooglePassed {
			t.Error("Google actually didn't pass")
		}
	} else {
		t.Error("Couldn't parse line")
	}
}

func TestHttp(t *testing.T) {
	//ts := startTestServer()
}
func startTestServer() *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "test_data/list_sample.txt")
	}))
	defer ts.Close()
	return ts
}
