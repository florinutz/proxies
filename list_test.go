package proxies

import (
	"bufio"
	"os"
	"testing"
)

func TestParser(t *testing.T) {
	f, err := os.Open("test_data/list_sample.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	parser := NewListParser(r)

	proxies := parser.GetProxies()

	if len(proxies) != 300 {
		t.Fatal("Not 300 proxies")
	}
}
