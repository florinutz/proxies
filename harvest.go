package proxies

import (
	"net/http"
)

type ListSourceReader interface {
	Populate(list *List) error
}

type ListWriter interface {
	Persist(list *List) error
}

type UrlSource *http.Response

type Source struct{}

func (s *Source) NewURLSource(URL string) (us UrlSource, err error) {
	res, err := http.Get(URL)
	us = UrlSource(res)
	if err != nil {
		return
	}
	if res.Status != "200 OK" {
		return
	}
	defer us.Body.Close()
	return
}

func (s *Source) Populate(l *List) error {
	return nil
}

func (s *Source) Persist(l *List) error {
	return nil
}
