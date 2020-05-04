package main

import (
	"fmt"
	"../utils"
)



func main() {
	// s := "hello"
	buf := make([]byte, 2)
	buf1 := append(buf, []byte("helloDY汉语")...)
	fmt.Println(string(buf1[2:len(buf1)]))

	bAddr := []byte("google.com:443")
	lenHeader := utils.Int8ToBytes(len(bAddr))
	
	// msg := append([]byte{len(bAddr)}, bAddr)
	fmt.Println(lenHeader)
	
}