package api

import (
	"encoding/json"
	"log"
	"testing"
)

func TestGetJMParam(t *testing.T) {
	var err error
	var jm JsonMap
	var b string
	var got any
	var tests []struct {
		path string
		want any
	}

	b = `
	{
		"apiVersion": "0.0",
		"Status": 200,
		"content": {
			"user": "tom@bobcat.com",
			"password": "very_secret_password",
			"new": {
				"password": "new_very_secret_password"
			}
		},
		"arr": [
			1,
			"some",
			true
			],
		"bool": true
	}	
	`
	tests = []struct {
		path string
		want any
	}{
		{path: "int", want: nil},
		{path: "Status", want: float64(200)},
		{path: "bool", want: true},
		{path: "arr.0", want: float64(1)},
		{path: "arr.1", want: "some"},
		{path: "arr.2", want: true},
		{path: "content.user", want: "tom@bobcat.com"},
		{path: "content.password", want: "very_secret_password"},
		{path: "content.new.password", want: "new_very_secret_password"},
	}
	if err = json.Unmarshal([]byte(b), &jm); err != nil {
		log.Fatal(err)
	}

	for _, tt := range tests {
		if got = jm.Get(tt.path); got != tt.want {
			t.Errorf("GetJMParam(%v) = %v, want %v", tt.path, got, tt.want)
		}
		if err != nil {
			t.Error(err.Error())
		}
	}
}

func TestGetJMParam_GetArray(t *testing.T) {
	var err error
	var jm JsonMap
	var b string
	var got any
	var tests []struct {
		index int
		want  any
	}

	b = `
	{
		"apiVersion": "0.0",
		"Status": 200,
		"content": {
			"user": "tom@bobcat.com",
			"password": "very_secret_password",
			"new": {
				"password": "new_very_secret_password"
			}
		},
		"arr": [
			1,
			"some",
			true
			],
		"bool": true
	}	
	`
	tests = []struct {
		index int
		want  any
	}{
		{index: 0, want: float64(1)},
		{index: 1, want: "some"},
		{index: 2, want: true},
	}
	if err = json.Unmarshal([]byte(b), &jm); err != nil {
		log.Fatal(err)
	}

	for _, tt := range tests {
		got = jm.Get("arr")
		if err != nil {
			t.Error(err.Error())
		}
		if (got.([]interface{}))[tt.index] != tt.want {
			t.Errorf("arr[%v] = %v, want %v", tt.index, got.([]interface{})[tt.index], tt.want)
		}
	}
}

func TestGetJMParam_ReplaceVlue(t *testing.T) {
	var err error
	var jm JsonMap
	var b string
	var got any
	var tests []struct {
		index   int
		replace any
	}

	b = `
	{
		"arr": [
			1,
			"some",
			true
			]
	}	
	`
	tests = []struct {
		index   int
		replace any
	}{
		{index: 0, replace: "some"},
		{index: 1, replace: 1},
		{index: 2, replace: false},
	}
	if err = json.Unmarshal([]byte(b), &jm); err != nil {
		log.Fatal(err)
	}

	for _, tt := range tests {
		got = jm.Get("arr")
		(got.([]interface{}))[tt.index] = tt.replace
		if err != nil {
			t.Error(err.Error())
		}
		if (got.([]interface{}))[tt.index] != (jm["arr"].([]interface{}))[tt.index] {
			t.Errorf("got[arr][%v] = %v, jm[arr][%v] = %v", tt.index, got.([]interface{})[tt.index], tt.index, (jm["arr"].([]interface{}))[tt.index])
		}
	}
}

func TestGetJMParam_TypeAssertion(t *testing.T) {
	var err error
	var jm JsonMap
	var b string
	var got any
	var tests []struct {
		index   string
		replace any
	}

	b = `
	{
		"obj": {
			"int": 1,
			"string": "some",
			"bool": true
			}
	}	
	`
	tests = []struct {
		index   string
		replace any
	}{
		{index: "int", replace: "some"},
		{index: "string", replace: 1},
		{index: "bool", replace: false},
	}
	if err = json.Unmarshal([]byte(b), &jm); err != nil {
		log.Fatal(err)
	}

	for _, tt := range tests {
		got = jm.Get("obj")
		(got.(map[string]interface{}))[tt.index] = tt.replace
		if err != nil {
			t.Error(err.Error())
		}
		if (got.(map[string]interface{}))[tt.index] != (jm["obj"].(map[string]interface{}))[tt.index] {
			t.Errorf("got[obj][%v] = %v, jm[obj][%v] = %v", tt.index, got.(map[string]interface{})[tt.index], tt.index, (jm["obj"].(map[string]interface{}))[tt.index])
		}
	}
}
