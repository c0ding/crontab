package main

import (
	"fmt"
	redis "github.com/garyburd/redigo/redis"
)

var (
	conn  redis.Conn
	err   error
	reply interface{}
)

func main() {

	if conn, err = redis.Dial("tcp", "localhost:6379"); err != nil {
		fmt.Println(err)

	}

	set()
	defer conn.Close()
}

func set() {
	//if reply, err = conn.Do("sadd", "set-a", "a", "b", "c"); err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(reply)

	//if reply, err = conn.Do("sadd", "set-b", "3", "b", "c", "d", "e"); err != nil {
	//	fmt.Println(err)
	//}

	if reply, err = conn.Do("sinter", "set-a", "set-b"); err != nil {
		fmt.Println("err ", err)
	}

	fmt.Println(redis.Strings(reply, err))
	strings, _ := redis.Strings(reply, err)
	for _, str := range strings {
		fmt.Println(str)
	}

}

func list() {
	//reply, err = conn.Do("lpush", "arr", "a")
	//fmt.Println(reply)

	if reply, err = redis.Strings(conn.Do("lrange", "arr", 0, -1)); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(reply)

}
func getAndSet() {
	reply, err = conn.Do("set", "name", "feifei")
	fmt.Println(reply)

	reply, err = redis.String(conn.Do("get", "name"))

	fmt.Printf("%v", reply)
}
