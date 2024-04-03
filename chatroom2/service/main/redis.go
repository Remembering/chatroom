package main

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

// 定义一个全局的pool
var pool *redis.Pool

func initPool(address string, maxIdle, maxActice int, idleTimeout time.Duration) {
	pool = &redis.Pool{
		MaxIdle:     maxIdle,     // 最大空闲链接数
		MaxActive:   maxActice,   //表示和数据库的最大链接数， 0 表示没有限制
		IdleTimeout: idleTimeout, //最大空闲时间,
		Dial: func() (redis.Conn, error) { //初始化链接的代码， 连接到redis
			return redis.Dial("tcp", address)
		},
	}
}
