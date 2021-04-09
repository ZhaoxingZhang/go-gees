package main

type RedisCli interface {
	Set(k, v interface{})
	Get(k interface{}) interface{}
	Del(k interface{})
}
