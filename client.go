package api

import (
	"encoding/json"
	"io"
	"net/http"
)

type IApi interface {
	Request(method string, data io.Reader) error
	Body() (ApiBody, error)
	CallApi(method string, data Request) (IApiBody, error)
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

func (c Client) Url() string {
	return c.url
}

func (c Client) Secret() string {
	return c.secret
}

func (c Client) UA() string {
	return c.userAgent
}

func (c Client) Version() string {
	return c.apiVersion
}

func (c Client) HttpResponse() *http.Response {
	return c.response
}

func (c *Client) Call(method string, data io.Reader) error {
	var err error
	var request *http.Request

	request, err = http.NewRequest(method, c.url, data)
	if err != nil {
		return err
	}
	if len(c.userAgent) > 0 {
		request.Header.Set("User-Agent", c.userAgent)
	}
	if len(c.secret) > 0 {
		request.Header.Set("X-Secret", c.secret)
	}
	c.response, err = http.DefaultClient.Do(request)
	return err
}

func (c *Client) makeCall(method string, data io.Reader) (ApiBody, error) {
	var err error
	var res ApiBody

	err = c.Call(method, data)
	if err != nil {
		return res, err
	}
	res, err = c.Body()
	if err != nil {
		return res, err
	}
	return res, err
}

func (c *Client) CallApi(method string, data Request) (ApiBody, error) {
	var err error
	var res ApiBody
	var reader io.Reader

	if reader, err = data.ToReader(); err != nil {
		return res, err
	}
	res, err = c.makeCall(method, reader)
	return res, err
}

func (c Client) Body() (ApiBody, error) {
	var data ApiBody
	var err error
	var dec *json.Decoder

	dec = json.NewDecoder(c.response.Body)
	err = dec.Decode(&data)
	return data, err
}
