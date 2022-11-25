package api

import (
	"bytes"
	"encoding/json"
	"io"
)

type Response struct {
	Status int     `json:"status"`
	Body   JsonMap `json:"body"`
}

func (r *Response) ToByteSlice() ([]byte, error) {
	var err error
	var res []byte

	res, err = json.Marshal(r)
	return res, err
}

func (r *Response) ToString() (string, error) {
	var err error
	var res []byte

	res, err = r.ToByteSlice()
	return string(res), err
}

func (r *Response) ToReader() (io.Reader, error) {
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
