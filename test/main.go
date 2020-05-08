package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strings"
)

func Int16ToBytes(n int) []byte {
	data := int16(n)
	bytebuf := bytes.NewBuffer([]byte{})
	binary.Write(bytebuf, binary.BigEndian, data)
	return bytebuf.Bytes()
}

func ioCopy(input net.Conn, output strings.Builder) (err error) {
	buf := make([]byte, 2)
	for {
		n, err := io.ReadAtLeast(input, buf, 2)
		if err != nil {
			fmt.Println(err)
			break
		}
		len := binary.BigEndian.Uint16(buf[0:2])
		if n >= 2 {
			fmt.Println("length:", len)
		}
		dataBuf := make([]byte, len)
		n, err = io.ReadFull(input, dataBuf)
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println(dataBuf)
	}
	return
}

func main() {
	var a = make([]byte, 10)
	var b = make([]byte, 10)
	fmt.Println(append(a, b...))

}
