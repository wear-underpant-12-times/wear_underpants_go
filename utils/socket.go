package utils

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

func PackHeader(data []byte) (msgData []byte, err error) {
	encoded := base64.StdEncoding.EncodeToString(data)
	encryptoS, err := AesEncrypt(encoded)
	if err != nil {
		return nil, err
	}
	encryptoB := []byte(encryptoS)
	dataLen := Int8ToBytes(len(encryptoB))
	msg := append(dataLen, encryptoB...)
	return msg, nil
}

func PackData(data []byte) (msgData []byte, err error) {
	encoded := base64.StdEncoding.EncodeToString(data)
	encryptoS, err := AesEncrypt(encoded)
	if err != nil {
		return nil, err
	}
	encryptoB := []byte(encryptoS)
	dataLen := Int16ToBytes(len(encryptoB))
	msg := append(dataLen, encryptoB...)
	return msg, nil
}

func UnPackData(data []byte) (msg []byte, err error) {
	decryptoS, err := AesDecrypt(string(data))
	if err != nil {
		return nil, err
	}
	decoded1, err := base64.StdEncoding.DecodeString(decryptoS)
	if err != nil {
		return nil, err
	}
	return decoded1, nil
}

func NetEncodeCopy(input net.Conn, output net.Conn) (err error) {
	buf := make([]byte, 8192)
	for {
		count, err := input.Read(buf)
		if err != nil {
			if err == io.EOF && count > 0 {
				data, err := PackData(buf[:count])
				if err != nil {
					continue
				}
				output.Write(data)
			}
			break
		}
		if count > 0 {
			data, err := PackData(buf[:count])
			if err != nil {
				continue
			}
			output.Write(data)
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
			upData, err := UnPackData(dataBuf)
			if err != nil {
				fmt.Println(err)
			} else {
				output.Write(upData)
			}
		}
	}
	return
}
