package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

func sendResponse(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	fmt.Fprintln(w, msg)
}

func apiError(w http.ResponseWriter, statusCode int, errStr string) {
	err := Error(statusCode, errStr)
	sendResponse(w, err.Error(), statusCode)
}

func ApiError(w http.ResponseWriter, statusCode int, errStr string) {
	apiError(w, statusCode, errStr)
}

func ServerFault(w http.ResponseWriter, err string) {
	var errStr = "internal server error"
	apiError(w, http.StatusInternalServerError, errStr)
}

func Forbidden(w http.ResponseWriter) {
	var err = "forbidden"
	defer apiError(w, http.StatusForbidden, err)
}

func BadRequest(w http.ResponseWriter) {
	var err = "bad request"
	apiError(w, http.StatusBadRequest, err)
}

func Unauthorized(w http.ResponseWriter) {
	var err = "unuthorized"
	apiError(w, http.StatusUnauthorized, err)
}

func NotFound(w http.ResponseWriter) {
	var err = "not found"
	apiError(w, http.StatusNotFound, err)
}

func NotAllowed(w http.ResponseWriter) {
	var err = "method not allowed"
	apiError(w, http.StatusMethodNotAllowed, err)
}

func Error(status int, msg string) error {
	var err error
	var body string
	var response = &Response{}

	response.Status = status
	response.Body = JsonMap{"msg": msg}
	if body, err = response.ToString(); err != nil {
		log.Println(err.Error())
	}

	return errors.New(body)
}
