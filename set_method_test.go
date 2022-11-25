package api

import (
	"encoding/json"
	"log"
	"testing"
)

func TestSetJMParam_SetArrVlue(t *testing.T) {
	var err error
	var jm JsonMap
	var got any
	var b string
	var tests []struct {
		key     string
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
		key     string
		replace any
	}{
		{key: "arr.0", replace: "some"},
		{key: "arr.1", replace: 1},
		{key: "arr.2", replace: false},
	}
	if err = json.Unmarshal([]byte(b), &jm); err != nil {
		log.Fatal(err)
	}

	for _, tt := range tests {
		err = jm.Set(tt.key, tt.replace)
		if err != nil {
			t.Error(err.Error())
		}
		got = jm.Get(tt.key)
		if got != tt.replace {
			t.Errorf("jm[%v] = %v, want: %v", tt.key, got, tt.replace)
		}
		if err != nil {
			t.Errorf("jm[%v] = %v, want: %v", tt.key, got, tt.replace)
		}
	}
}

func TestSetJMParam_SetMapVlue(t *testing.T) {
	var err error
	var jm JsonMap
	var b string
	var got any
	var tests []struct {
		key     string
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
		key     string
		replace any
	}{
		{key: "obj.int", replace: "some"},
		{key: "obj.string", replace: 1},
		{key: "obj.bool", replace: false},
	}
	if err = json.Unmarshal([]byte(b), &jm); err != nil {
		log.Fatal(err)
	}

	for _, tt := range tests {
		err = jm.Set(tt.key, tt.replace)
		if err != nil {
			t.Error(err.Error())
		}
		got = jm.Get(tt.key)
		if got != tt.replace {
			t.Errorf("jm[%v] = %v, want: %v", tt.key, got, tt.replace)
		}
		if err != nil {
			t.Errorf("jm[%v] = %v, want: %v", tt.key, got, tt.replace)
		}
	}
}

func TestSetJMParam(t *testing.T) {
	var err error
	var jm JsonMap
	var b string
	var got any
	var tests []struct {
		key     string
		set     any
		replace any
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
		key     string
		set     any
		replace any
	}{
		{key: "Status", replace: float64(404)},
		{key: "bool", replace: false},
		{key: "arr.0", replace: float64(15)},
		{key: "arr.1", replace: "some string"},
		{key: "arr.2", replace: false},
		{key: "content.user", replace: "bob@bobcat.com"},
		{key: "content.password", replace: "not_so_secret"},
		{key: "content.new.password", replace: "not_very_secret_password"},
	}
	if err = json.Unmarshal([]byte(b), &jm); err != nil {
		log.Fatal(err)
	}

	for _, tt := range tests {
		err = jm.Set(tt.key, tt.replace)
		if err != nil {
			t.Error(err.Error())
		}
		got = jm.Get(tt.key)
		if got != tt.replace {
			t.Errorf("jm[%v] = %v, want: %v", tt.key, got, tt.replace)
		}
		if err != nil {
			t.Errorf("jm[%v] = %v, want: %v", tt.key, got, tt.replace)
		}
	}
}
