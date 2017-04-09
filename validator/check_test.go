package validator

import (
	"github.com/elazarl/goproxy"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

const Body = "sup"

type SimpleStringHanlder string

func (h SimpleStringHanlder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, string(h))
}

func TestCheck(t *testing.T) {
	proxyServer := httptest.NewServer(goproxy.NewProxyHttpServer())
	defer proxyServer.Close()

	var httpServer = httptest.NewServer(SimpleStringHanlder(Body))
	check(httpServer, proxyServer, t)

	var httpsServer = httptest.NewTLSServer(SimpleStringHanlder(Body))
	check(httpsServer, proxyServer, t)
}

func check(httpServer *httptest.Server, proxyServer *httptest.Server, t *testing.T) {
	resp, err := Check(proxyServer.URL, httpServer.URL)
	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if string(body) != Body {
		t.Fatalf("Expected '%s', got '%s'", Body, string(body))
	}
}
