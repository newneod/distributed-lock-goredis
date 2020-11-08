# distributed-lock
The distributed-lock is a simple distributed lock.This version was developed by go-redis.
You can find locks developed with redigo in my other projects.

## download and install
```
$ go get github.com/newneod/distributed-lock-goredis
```

## demo
```
package main

import "time"

// demo
func main() {
	Init("127.0.0.1:6379")
	defer conn.Close()

	strUUID, err := Lock("a")
	if err != nil {
		panic(err)
	}

	time.Sleep(3 * time.Second)
	err = Unlock("a", string(strUUID))
	if err != nil {
		panic(err)
	}
}
```