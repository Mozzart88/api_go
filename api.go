package api

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
)

type IApiBody interface {
	ToByteSlice() []byte
	ToString() string
	ToReader() io.Reader
}

type ApiBody struct {
	status int
	body   JsonMap
}

func (r ApiBody) ToByteSlice() ([]byte, error) {
	var err error
	var res []byte

	res, err = json.Marshal(r)
	return res, err
}

func (r ApiBody) ToString() (string, error) {
	var err error
	var res []byte

	res, err = r.ToByteSlice()
	return string(res), err
}

func (r ApiBody) ToReader() (io.Reader, error) {
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

func splitPath(path string) (string, string) {
	var begin string
	var end string
	var sep string = "."
	var keys []string

	keys = strings.Split(path, sep)
	if len(keys) == 1 {
		return begin, path
	}
	begin = strings.Join(keys[0:len(keys)-1], sep)
	end = keys[len(keys)-1]
	return begin, end
}
