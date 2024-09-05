package prugo

import (
	"net/http"
	"slices"
	"time"
)

type Config struct {
	headers http.Header
}

type Prugo interface {
	Get(string, Config) (*http.Response, error)
	Post()
	Put()
	Patch()
	Delete()
}

type prugo struct {
	baseURL string
	headers http.Header
	timeout time.Duration
}

func New(baseURL string, headers http.Header, timeout time.Duration) Prugo {

	p := &prugo{

		timeout: time.Second,
		headers: http.Header{
			"Content-Type": {"application/json"},
		},
	}

	if baseURL != "" {
		p.baseURL = baseURL
	}

	if headers != nil {
		p.headers = headers
	}

	if timeout != time.Second*0 {
		p.timeout = timeout
	}

	return p
}

func (p *prugo) Get(url string, config Config) (*http.Response, error) {

	requestURL := p.baseURL + url

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}

	if config.headers != nil {
		for k, v := range p.headers {
			_, ok := config.headers[k]
			if !ok {
				config.headers[k] = v
			} else {
				slices.Concat(config.headers[k], v)
			}
		}
	}

	req.Header = config.headers

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil

}

func (p *prugo) Post() {

}

func (p *prugo) Put() {

}

func (p *prugo) Patch() {

}

func (p *prugo) Delete() {

}
