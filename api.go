package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type JsonMap map[string]interface{}

type ApiResponse struct {
	Status int
	Body   map[string]any
}

type Server struct {
	Url       string
	Secret    string
	UserAgent string
}

func ApiError(w http.ResponseWriter, statusCode int, content string) {
	var response ApiResponse
	var data []byte
	var err error
	var body = map[string]any{
		"msg": content,
	}

	response = ApiResponse{Body: body, Status: statusCode}
	data, err = json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Error(w, string(data), statusCode)
}

func (s *Server) MakeRequest(method string, data io.Reader) (*http.Response, error) {
	var err error
	var response *http.Response
	var request *http.Request

	request, err = http.NewRequest(method, s.Url, data)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", s.UserAgent)
	response, err = http.DefaultClient.Do(request)
	return response, err
}

func ReadResponseBody(body io.ReadCloser) (ApiResponse, error) {
	var data ApiResponse
	var err error
	var dec *json.Decoder

	dec = json.NewDecoder(body)
	err = dec.Decode(&data)
	return data, err
}

func Map2jsonByte(data JsonMap) ([]byte, error) {
	var res []byte
	var err error

	res, err = json.Marshal(data)
	return res, err
}

func Map2jsonReader(data JsonMap) (io.Reader, error) {
	var jsonRes []byte
	var err error
	var res io.Reader

	jsonRes, err = Map2jsonByte(data)
	if err == nil {
		res = bytes.NewReader(jsonRes)
	}
	return res, err
}

func (api ApiResponse) ToByteSlice() []byte {
	var err error
	var res []byte

	res, err = json.Marshal(api)
	if err != nil {
		log.Fatalf("Fail to Marshal api struct: %s", err.Error())
		return res
	}
	return res
}

func (api JsonMap) ToByteSlice() []byte {
	var err error
	var res []byte

	res, err = json.Marshal(api)
	if err != nil {
		log.Fatalf("Fail to Marshal api struct: %s", err.Error())
		return res
	}
	return res
}

func (s Server) GetRequestBody(body io.Reader) (JsonMap, error) {
	var err error
	var data JsonMap
	var b []byte

	b, err = io.ReadAll(body)
	if err != nil {
		return data, err
	}
	if !json.Valid(b) {
		return data, errors.New("invalid json")
	}
	err = json.Unmarshal(b, &data)
	return data, err
}

func (s *Server) VerifyUA(accept []string) bool {
	var v string

	for _, v = range accept {
		if v == s.UserAgent {
			return true
		}
	}
	return false
}

func (s *Server) callServer(method string, data io.Reader) (ApiResponse, error) {
	var err error
	var response *http.Response
	var res ApiResponse

	response, err = s.MakeRequest(method, data)
	if err != nil {
		return res, err
	}
	res, err = ReadResponseBody(response.Body)
	if err != nil {
		return res, err
	}
	return res, err
}

func (s *Server) CallApi(method string, data JsonMap) (ApiResponse, error) {
	var err error
	var res ApiResponse
	var reader io.Reader

	reader, err = Map2jsonReader(data)
	if err != nil {
		return res, err
	}
	res, err = s.callServer(method, reader)
	return res, err
}

func getValue(jm any, key any) (any, error) {
	var err error
	var val any

	switch jm.(type) {
	case []any:
		if key, err = strconv.Atoi(key.(string)); err != nil {
			return val, err
		}
		jm = (jm.([]any))[key.(int)]
	case map[string]any:
		jm = (jm.(map[string]any))[key.(string)]
	default:
		jm = (jm.(JsonMap))[key.(string)]
	}
	return jm, err
}

func (jm JsonMap) Get(path string) (any, error) {
	const sep string = "."
	var err error
	var val any
	var keys []string
	var key string
	var tmp any

	keys = strings.Split(path, sep)
	tmp = jm
	for _, key = range keys {
		tmp, err = getValue(tmp, key)
		if err != nil {
			return val, err
		}
		if tmp == nil {
			break
		}
	}
	val = tmp
	return val, err
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

func (jm JsonMap) Set(path string, value any) error {
	var err error
	var tmp any
	var key string

	path, key = splitPath(path)
	if len(path) == 0 {
		jm[key] = value
	} else {
		tmp, err = jm.Get(path)
		if err != nil {
			return err
		}
		if tmp == nil {
			return errors.New("Invalid index")
		}
		switch tmp.(type) {
		case []any:
			if err = setByIndex(tmp.([]any), key, value); err != nil {
				return err
			}
		case map[string]any:
			if err = setByKey(tmp.(map[string]any), key, value); err != nil {
				return err
			}
		default:
			jm[path] = value
		}
	}

	return err
}

func Error(status int, msg string) ApiResponse {
	var res = ApiResponse{
		Status: status,
		Body: map[string]any{
			"msg": msg,
		},
	}
	return res
}
