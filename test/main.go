package main

import (
	"fmt"
	"io"
	"strings"
	"bytes"
	"encoding/binary"
)

func Int16ToBytes(n int) []byte {
    data := int16(n)
    bytebuf := bytes.NewBuffer([]byte{})
    binary.Write(bytebuf, binary.BigEndian, data)
    return bytebuf.Bytes()
}


func main() {
	
	r := strings.NewReader("hello")
	buf := make([]byte, 14)

	io.ReadAtLeast(r, buf, 4)
	fmt.Println(buf)
}