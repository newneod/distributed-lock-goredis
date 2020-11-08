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
