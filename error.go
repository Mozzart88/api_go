package api

import (
	"bytes"
	"fmt"
	"io"
)

type Error struct {
	status int
	msg    string
}

func NewError(status int, msg string) error {
	return Error{status, msg}
}

func (e Error) ToString() string {
	var res string
	var format string = `{"status":%d,"body":{"msg":%s}}`

	res = fmt.Sprintf(format, e.status, e.msg)
	return res
}

func (e Error) Error() string {
	return e.ToString()
}

func (e Error) ToByteSlice() []byte {
	var res []byte

	res = []byte(e.ToString())
	return res
}

func (e Error) ToReader() io.Reader {
	var tmp []byte
	var res io.Reader

	tmp = e.ToByteSlice()
	res = bytes.NewReader(tmp)
	return res
}

func (e Error) Message() string {
	return e.msg
}

func (e Error) Status() int {
	return e.status
}
