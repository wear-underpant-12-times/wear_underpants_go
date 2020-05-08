package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"../utils"
)

var (
	Commands = []string{"CONNECT", "BIND", "UDP ASSOCIATE"}
	AddrType = []string{"", "IPv4", "", "Domain", "IPv6"}
	Conns    = make([]net.Conn, 0)
	Verbose  = false

	errAddrType      = errors.New("socks addr type not supported")
	errVer           = errors.New("socks version not supported")
	errMethod        = errors.New("socks only support noauth method")
	errAuthExtraData = errors.New("socks authentication get extra data")
	errReqExtraData  = errors.New("socks request get extra data")
	errCmd           = errors.New("socks only support connect command")
)

func shake(conn net.Conn) (target string, err error) {
	buf := make([]byte, 258)
	var n int
	if n, err = io.ReadAtLeast(conn, buf, 1); err != nil {
		fmt.Println(err)
		return
	}

	dmLen := int(buf[0])
	msgLen := dmLen + 1
	// fmt.Println(msgLen)
	if n == msgLen {
	} else if n < msgLen {
		if _, err = io.ReadFull(conn, buf[n:msgLen]); err != nil {
			fmt.Println("get full failed")
			return
		}
	} else {
		fmt.Printf("dmLen %v, getLen %v\n", dmLen, n)
		// return errors.New("auth error")
	}
	return string(buf[1 : 1+dmLen]), nil
}

func netCopy(input, output net.Conn) (err error) {
	buf := make([]byte, 8192)
	for {
		count, err := input.Read(buf)
		if err != nil {
			if err == io.EOF && count > 0 {
				output.Write(buf[:count])
			}
			break
		}
		if count > 0 {
			output.Write(buf[:count])
		}
	}
	return
}

func handConn(conn net.Conn) {
	defer conn.Close()
	addr, err := shake(conn)
	if err != nil {
		fmt.Println("shake error", err)
		return
	}
	fmt.Println(addr)
	remoteConn, err := net.DialTimeout("tcp", addr, time.Duration(time.Second*15))
	if err != nil {
		fmt.Println("connect server error:", addr, err)
		return
	}
	defer remoteConn.Close()
	go utils.NetDecodeCopy(conn, remoteConn)
	utils.NetEncodeCopy(remoteConn, conn)
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:8082")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handConn(conn)
	}
}
