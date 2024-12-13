package model

import (
	"github.com/astaxie/goredis"
)


var redisCli goredis.Client

func init() {
	redisCli.Addr = "127.0.0.1:6379"
}

func PutinQueue(key string, value string) error {
	return redisCli.Lpush(key, []byte(value))
}

func PopFromQueue(key string) string {
	res, err := redisCli.Rpop(key)
	if err != nil {
		panic(err)
	}
	return string(res)
}

func AddToSet(key string, value string) (bool, error) {
	return redisCli.Sadd(key, []byte(value))
}

func ISVist(key, value string) bool {
	// 判断是否在redis set里存储了
	set, err := Smembers(key)
	if err != nil{
		return false
	} else {
		for i := 0; i < len(set); i++ {
			if (string(set[i]) == value) {return true}
		}
	}
	return false
}

func Smembers(key string) ([][]byte, error){
	return redisCli.Smembers(key)
}
