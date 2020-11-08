package main

import (
	"github.com/go-redis/redis/v8"
)

var conn *redis.Client

func Init(address string) {
	if address == "" {
		address = "127.0.0.1:6379"
	}

	conn = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "",
		DB:       0,
	})
}
