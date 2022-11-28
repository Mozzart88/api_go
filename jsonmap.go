package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"strconv"
	"strings"
)

type JsonMap map[string]any

func (jm *JsonMap) ToByteSlice() ([]byte, error) {
	var err error
	var res []byte

	res, err = json.Marshal(jm)
	return res, err
}

func (jm *JsonMap) ToString() (string, error) {
	var err error
	var res []byte

	res, err = jm.ToByteSlice()
	return string(res), err
}

func (jm JsonMap) ToReader() (io.Reader, error) {
	var err error
	var tmp []byte
	var res io.Reader

	tmp, err = jm.ToByteSlice()
	if err != nil {
		return res, err
	}
	res = bytes.NewReader(tmp)
	return res, nil
}

func (jm JsonMap) Keys() []string {
	var keys []string

	for k := range jm {
		keys = append(keys, k)
	}
	return keys
}

func getValue(jm any, key any) any {
	var err error

	switch t := jm.(type) {
	case []any:
		if key, err = strconv.Atoi(key.(string)); err != nil {
			return err
		}
		jm = t[key.(int)]
	case map[string]any:
		jm = t[key.(string)]
	case JsonMap:
		jm = t[key.(string)]
	default:
		jm = t
	}
	return jm
}

/*
May return error - need to check if returned value is error
*/
func (jm JsonMap) Get(path string) any {
	const sep string = "."
	var err error
	var val any
	var keys []string
	var key string
	var tmp any

	keys = strings.Split(path, sep)
	tmp = jm
	for _, key = range keys {
		tmp = getValue(tmp, key)
		if tmp == nil {
			break
		}
		if _, ok := tmp.(error); ok {
			return err
		}
	}
	val = tmp
	return val
}

func setByIndex(sl []any, index string, value any) error {
	var err error
	var i int

	if i, err = strconv.Atoi(index); err != nil {
		return err
	}
	sl[i] = value
	return err
}

func setByKey(sl map[string]any, index string, value any) error {
	var err error

	sl[index] = value
	return err
}

/*
TODO: If as value parameter passed stringified json
we need to parse it and set values as json value and
not as string
*/
func (jm *JsonMap) Set(path string, value any) error {
	var err error
	var tmp any
	var key string

	path, key = splitPath(path)
	if len(path) == 0 {
		(*jm)[key] = value
	} else {
		tmp = jm.Get(path)
		if e, ok := tmp.(error); ok {
			return e
		}
		if tmp == nil {
			return errors.New("invalid index")
		}
		switch tmp := tmp.(type) {
		case []any:
			if err = setByIndex(tmp, key, value); err != nil {
				return err
			}
		case map[string]any:
			if err = setByKey(tmp, key, value); err != nil {
				return err
			}
		default:
			(*jm)[path] = value
		}
	}

	return err
}
