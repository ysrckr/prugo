package prugo

import (
	"net/url"
	"time"
)

type Prugo interface {
	Get()
	Post()
	Put()
	Patch()
	Delete()
}

type prugo struct {
	baseURL *url.URL
	headers map[string]string
	timeout time.Duration
}

func New(baseURL string, headers map[string]string, timeout time.Duration) (Prugo, error) {
	var err error

	p := &prugo{
		baseURL: &url.URL{},
		timeout: time.Second,
		headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	if baseURL != "" {
		p.baseURL, err = url.Parse(baseURL)
		if err != nil {
			return p, err
		}
	}

	if headers != nil {
		p.headers = headers
	}

	if timeout != time.Second*0 {
		p.timeout = timeout
	}

	return p, nil
}

func (p *prugo) Get() {

}

func (p *prugo) Post() {

}

func (p *prugo) Put() {

}

func (p *prugo) Patch() {

}

func (p *prugo) Delete() {

}
