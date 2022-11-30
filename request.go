package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type Request struct {
	body   JsonMap
	header map[string][]string
}

func newRequest(r *http.Request) (*Request, error) {
	var request = new(Request)
	var err error
	var b []byte
	var body JsonMap

	b, err = io.ReadAll(r.Body)
	if err == nil && len(b) > 0 {
		err = json.Unmarshal(b, &body)
	}
	request.body = body
	request.header = r.Header
	return request, err

}

func (r Request) Header(h string) []string {
	return r.header[h]
}

func (r Request) Body(path string) any {
	var res any

	// if len(path) == 0 {
	// 	return r.body
	// }
	res = r.body.Get(path)
	if _, ok := res.(error); ok {
		return nil
	}
	return res
}

func (r Request) IsEmptyBody() bool {
	return r.body == nil
}

func (r *Request) ToByteSlice() ([]byte, error) {
	var err error
	var res []byte

	res, err = json.Marshal(r.body)
	return res, err
}

func (r *Request) ToString() (string, error) {
	var err error
	var res []byte

	res, err = r.ToByteSlice()
	return string(res), err
}

func (r *Request) ToReader() (io.Reader, error) {
	var err error
	var tmp []byte
	var res io.Reader

	tmp, err = r.ToByteSlice()
	if err != nil {
		return res, err
	}
	res = bytes.NewReader(tmp)
	return res, nil
}
