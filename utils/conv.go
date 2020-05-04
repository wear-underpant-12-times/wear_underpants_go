package utils

import (
	"bytes"
	"encoding/binary"
)


func Int8ToBytes(n int) []byte {
    data := int8(n)
    bytebuf := bytes.NewBuffer([]byte{})
    binary.Write(bytebuf, binary.BigEndian, data)
    return bytebuf.Bytes()
}