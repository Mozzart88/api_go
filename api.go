package api

import (
	"io"
	"strings"
)

type ApiBody interface {
	ToByteSlice() []byte
	ToString() string
	ToReader() io.Reader
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
