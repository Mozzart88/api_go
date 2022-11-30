package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func emptyBody(w http.ResponseWriter, r *http.Request) {
	srv := NewServer(w, r)

	if srv.Request.IsEmptyBody() {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("KO"))
	}
}

func GetHoleBody(w http.ResponseWriter, r *http.Request) {
	srv := NewServer(w, r)
	body := srv.Request.Body("")
	var ok bool
	var b JsonMap

	if b, ok = body.(JsonMap); !ok {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Expected type JsonMap, got %T", body)))
	} else {
		bb, _ := io.ReadAll(r.Body)
		bi, _ := b.ToByteSlice()
		if bytes.Equal(bb, bi) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("%v !+= %v", bb, bi)))
		}
	}

}

func TestEmptyBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/empty/", nil)
	w := httptest.NewRecorder()
	emptyBody(w, req)
	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Error("expected empty body")
	}
}

func TestNotEmptyBody(t *testing.T) {
	body := bytes.NewReader([]byte(`{"msg":"ok"}`))
	req := httptest.NewRequest(http.MethodGet, "/empty/", body)
	w := httptest.NewRecorder()
	emptyBody(w, req)
	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
		t.Error("expected not empty body")
	}
}

func TestGetEmptyPath(t *testing.T) {
	body := bytes.NewReader([]byte(`{"msg":"ok"}`))
	req := httptest.NewRequest(http.MethodGet, "/empty/", body)
	w := httptest.NewRecorder()
	emptyBody(w, req)
	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
		b, _ := io.ReadAll(res.Body)
		t.Error(string(b))
	}
}
