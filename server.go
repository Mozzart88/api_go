package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Server struct {
	responseWriter http.ResponseWriter
	request        *http.Request
}

func NewServer(w http.ResponseWriter, r *http.Request) Server {
	return Server{w, r}
}

func (s Server) Body() (JsonMap, error) {
	var err error
	var data JsonMap
	var b []byte

	b, err = io.ReadAll(s.request.Body)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(b, &data)
	return data, err
}

func (s Server) UnmarshalBody(v any) error {
	var err error
	var b []byte

	b, err = io.ReadAll(s.request.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, v)
	return err
}

func (s *Server) VerifyUA(accept []string) bool {
	var v string

	for _, v = range accept {
		if v == s.request.UserAgent() {
			return true
		}
	}
	return false
}

func (s Server) Forbidden() {
	Forbidden(s.responseWriter)
	log.Println(s.request)
}

func (s Server) ServerFault(err error) {
	ServerFault(s.responseWriter, err.Error())
	log.Println(s.request)
}

func (s Server) BadRequest() {
	BadRequest(s.responseWriter)
	log.Println(s.request)
}

func (s Server) NotAllowed() {
	NotAllowed(s.responseWriter)
	log.Println(s.request)
}

func (s Server) Error(status int, msg string) {
	apiError(s.responseWriter, status, msg)
	log.Println(s.request)
}

func (s Server) Responed(status int, msg string) {
	sendResponse(s.responseWriter, msg, status)
}

func (s Server) Method() string {
	return s.request.Method
}
