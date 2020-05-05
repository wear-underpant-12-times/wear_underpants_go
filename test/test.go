package main

import (
	"fmt"
	"strconv"
)

func main() {
	var buf = []byte{1, 3}
	fmt.Println(strconv.Itoa(int(buf[1])))
}
