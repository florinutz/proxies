package proxies

import (
	"bufio"
	"os"
	"testing"
)

func TestReadFile(t *testing.T) {
	list := List{}

	f, err := os.Open("test_data/list_sample.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	list.Read(r)

	// there are 300 possibly valid proxiesRead in list_sample.txt
	if len(list) != 300 {
		t.Fatal("Not 300 proxiesRead")
	}
}

func TestReadLine(t *testing.T) {
	l := List{}
	if proxy := l.ReadSingleProxy("54.153.73.116:3128 US-N-S! -"); proxy != nil {
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
