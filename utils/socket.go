package utils

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

func NetEncodeCopy(input net.Conn, output net.Conn) (err error) {
	buf := make([]byte, 8192)
	for {
		count, err := input.Read(buf)
		if err != nil {
			if err == io.EOF && count > 0 {
				header := Int16ToBytes(count)
				msg := append(header, buf[:count]...)
				output.Write(msg)
			}
			break
		}
		header := Int16ToBytes(count)
		msg := append(header, buf[:count]...)
		if count > 0 {
			output.Write(msg)
		}
	}
	return
}

func NetDecodeCopy(input net.Conn, output net.Conn) (err error) {
	buf := make([]byte, 2)
	for {
		n, err := io.ReadAtLeast(input, buf, 2)
		if err != nil {
			fmt.Println(err)
			break
		}
		if n < 2 {
			break
		}
		dataLen := binary.BigEndian.Uint16(buf)
		dataBuf := make([]byte, dataLen)
		n, err = io.ReadFull(input, dataBuf)
		if err != nil {
			fmt.Println(err)
			break
		}
		if n > 0 {
			output.Write(dataBuf)
		}
	}
	return
}
