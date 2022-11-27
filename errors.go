package api

import (
	"fmt"
	"log"
	"net/http"
)

func (e Error) Send(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(e.status)
	fmt.Fprintln(w, e.ToString())
}

func (e Error) Log() {
	log.Println(e.Message())
}

func (e Error) SendAndLog(w http.ResponseWriter) {
	e.Log()
	e.Send(w)
}

func apiError(statusCode int, errStr string) Error {
	var err Error = NewError(statusCode, errStr).(Error)
	return err
}

func appendErrMsg(start string, chunks []string) string {
	var chunk string
	var i, l int

	if l = len(chunks) - 1; l >= 0 {
		start += ": "
	}
	for i, chunk = range chunks {
		start += chunk
		if i < l {
			start += ", "
		}
	}
	return start
}

func ApiError(statusCode int, errStr string) Error {
	return apiError(statusCode, errStr)
}

func ServerFault(msgs ...string) Error {
	var err = "internal server error"

	err = appendErrMsg(err, msgs)
	return apiError(http.StatusInternalServerError, err)
}

func Forbidden(msgs ...string) Error {
	var err = "forbidden"
	err = appendErrMsg(err, msgs)
	return apiError(http.StatusForbidden, err)
}

func BadRequest(msgs ...string) Error {
	var err = "bad request"
	err = appendErrMsg(err, msgs)
	return apiError(http.StatusBadRequest, err)
}

func Unauthorized(msgs ...string) Error {
	var err = "unuthorized"
	err = appendErrMsg(err, msgs)
	return apiError(http.StatusUnauthorized, err)
}

func NotFound(msgs ...string) Error {
	var err = "not found"
	err = appendErrMsg(err, msgs)
	return apiError(http.StatusNotFound, err)
}

func NotAllowed(msgs ...string) Error {
	var err = "method not allowed"
	err = appendErrMsg(err, msgs)
	return apiError(http.StatusMethodNotAllowed, err)
}
