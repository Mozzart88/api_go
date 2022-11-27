package api

import (
	"io"
	"net/http"
	"strings"
)

type IApi interface {
	Request(method string, data io.Reader) error
	Body() (Response, error)
	CallApi(method string, data Request) (Response, error)
}

type Client struct {
	url        string
	secret     string
	userAgent  string
	apiVersion string
	response   *http.Response
}

func NewClient(url string, apiVersion string) *Client {
	return &Client{url: url, apiVersion: apiVersion}
}

func (c Client) Url(url ...string) string {
	if len(url) > 0 {
		c.url = url[0]
	}
	return c.url
}

func (c Client) Secret(sec ...string) string {
	if len(sec) > 0 {
		c.secret = sec[0]
	}
	return c.secret
}

func (c Client) UA(ua ...string) string {
	if len(ua) > 0 {
		c.userAgent = ua[0]
	}
	return c.userAgent
}

func (c Client) Version(v ...string) string {
	if len(v) > 0 {
		c.apiVersion = v[0]
	}
	return c.apiVersion
}

func (c Client) HttpResponse() *http.Response {
	return c.response
}

func (c Client) Call(method string, path string, data io.Reader) (*http.Response, error) {
	var err error
	var request *http.Request
	var url string = strings.Join([]string{c.url, path}, "/")

	request, err = http.NewRequest(method, url, data)
	if err != nil {
		return nil, err
	}
	if len(c.userAgent) > 0 {
		request.Header.Set("User-Agent", c.userAgent)
	}
	if len(c.secret) > 0 {
		request.Header.Set("X-Secret", c.secret)
	}
	return http.DefaultClient.Do(request)
}

func (c Client) makeCall(method string, path string, data io.Reader) (*http.Response, error) {
	return c.Call(method, path, data)
}

func (c Client) CallApi(method string, path string, data JsonMap) (*http.Response, error) {
	var err error
	var reader io.Reader

	if reader, err = data.ToReader(); err != nil {
		return nil, err
	}
	return c.makeCall(method, path, reader)
}
