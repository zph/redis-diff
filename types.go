package main

import (
	"gopkg.in/redis.v4"
)

type redisKey string
type pattern string

type Pipe struct {
	from     *Server
	to       *Server
	keys     string
	shutdown chan bool
}

type Server struct {
	client *redis.Client
	host   string
	port   int
	db     int
	pass   string
}

type Discrepancy struct {
	key redisKey
	src interface{}
	dst interface{}
}
