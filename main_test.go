package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func assertFieldsEqual(t *testing.T, e *Server, a *Server) {
	assert.Equal(t, e.host, a.host, "host should be set")
	assert.Equal(t, e.port, a.port, "port should be set")
	assert.Equal(t, e.db, a.db, "db should be set")
	assert.Equal(t, e.pass, a.pass, "password should be set")
}

func TestParseURI(t *testing.T) {
	s := "redis://x:password@host.com:123"
	a, _ := parseURI(s)
	e := &Server{
		host: "host.com",
		port: 123,
		db:   0,
		pass: "password",
	}
	assertFieldsEqual(t, e, a)

	s3 := "redis://localhost:6370"
	a3, _ := parseURI(s3)
	e3 := &Server{nil, "localhost", 6370, 0, ""}
	assertFieldsEqual(t, e3, a3)
}

func TestRedisURIParsing(t *testing.T) {
	s := "redis://x:password@host.com:123"
	a, _ := parseURI(s)
	e := &Server{
		host: "host.com",
		port: 123,
		db:   0,
		pass: "password",
	}
	assertFieldsEqual(t, e, a)
}

func TestRedisURIParsingWithDB(t *testing.T) {
	s := "redis://x:password@host.com:123?db=0"
	actual, _ := parseURI(s)
	expected := &Server{
		host: "host.com",
		port: 123,
		db:   0,
		pass: "password",
	}
	assertFieldsEqual(t, expected, actual)
}
