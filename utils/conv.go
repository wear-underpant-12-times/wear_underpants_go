package utils

import (
	"bytes"
	"encoding/binary"
	"errors"
)

func Int8ToBytes(n int) []byte {
	data := int8(n)
	bytebuf := bytes.NewBuffer([]byte{})
	binary.Write(bytebuf, binary.BigEndian, data)
	return bytebuf.Bytes()
}

func Int16ToBytes(n int) []byte {
	data := int16(n)
	bytebuf := bytes.NewBuffer([]byte{})
	binary.Write(bytebuf, binary.BigEndian, data)
	return bytebuf.Bytes()
}

func MergeBytes(list [][]byte) ([]byte, error) {
	if len(list) < 1 {
		return nil, errors.New("len error")
	}
	res := []byte{}
	for _, bs := range list {
		res = append(res, bs...)
	}
	return res, nil
}
