package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strings"

	"../utils"
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
	s, err := utils.AesEncrypt("4天前 - 读者在线阅读,读者文摘在线阅读,《读者》其栏目文苑,人物,名人轶事,社会,人生,生活点滴更赢得了读者的喜爱与拥护,因此《读者》被誉为中国人的心灵读本...")
	if err != nil {
		panic("err")
	}
	fmt.Println(s)
	d, err := utils.AesDecrypt(s)
	if err != nil {
		panic("de error")
	}
	fmt.Println(d)
}
