package api

import (
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	responseWriter http.ResponseWriter
	request        *http.Request
	Request        *Request
}

func NewServer(w http.ResponseWriter, r *http.Request) Server {
	var srv Server
	var err error

	srv.responseWriter = w
	srv.request = r
	if srv.Request, err = newRequest(r); err != nil {
		srv.ServerFault(err.Error())
	}
	return srv
}

func (s Server) Responder() http.ResponseWriter {
	return s.responseWriter
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

func (s Server) responedParseResp(response []any) (string, map[string]string) {
	var msg string
	var headers map[string]string

	for _, r := range response {
		switch t := r.(type) {
		case string:
			msg = t
		case map[string]string:
			headers = t
		default:
			panic("only string for message and map[string]string for headers allowed")
		}
	}
	return msg, headers
}

func (s Server) Forbidden(response ...any) {
	var msg string
	var headers map[string]string

	msg, headers = s.responedParseResp(response)
	s.Error(http.StatusForbidden, msg, headers)
}

func (s Server) ServerFault(response ...any) {
	var msg string
	var headers map[string]string

	msg, headers = s.responedParseResp(response)
	s.Error(http.StatusInternalServerError, msg, headers)
}

func (s Server) BadRequest(response ...any) {
	var msg string
	var headers map[string]string

	msg, headers = s.responedParseResp(response)
	s.Error(http.StatusBadRequest, msg, headers)
}

func (s Server) NotAllowed(response ...any) {
	var msg string
	var headers map[string]string

	msg, headers = s.responedParseResp(response)
	s.Error(http.StatusMethodNotAllowed, msg, headers)
}

func (s Server) Error(status int, msg string, headers ...map[string]string) {
	if len(headers) > 0 {
		s.Responed(status, msg, headers[0])
	} else {
		s.Responed(status, msg)
	}

	log.Println(s.request)
}

func (s Server) Responed(status int, msg string, headers ...map[string]string) {
	var w = s.responseWriter
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	if len(headers) > 0 {
		for _, m := range headers {
			for h, v := range m {
				w.Header().Set(h, v)
			}
		}
	}
	w.WriteHeader(status)
	fmt.Fprintln(w, msg)
}

func (s Server) Method() string {
	return s.request.Method
}
