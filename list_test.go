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

	parser := NewListParser(bufio.NewReader(f))

	proxies := parser.GetProxies()

	// there are 300 possibly valid proxies in list_sample.txt
	if len(proxies) != 300 {
		t.Fatal("Not 300 proxies")
	}
}
