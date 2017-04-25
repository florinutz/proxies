package proxies

import (
	"io"
	"net/http"
)

type Reader interface {
	Read(reader io.Reader, list *List)
}

type HttpClient struct {
	*http.Client
	List   *List
	Reader Reader
}

func (c *HttpClient) Do(req *http.Request) (*http.Response, error) {
	res, err := c.Client.Do(req)
	if err == nil {
		c.Reader.Read(res.Body, c.List)
	}
	return res, err
}

type LineReader func(reader io.Reader, list *List)

func (lr LineReader) Read(reader io.Reader, list *List) {
	list.Read(reader)
}
