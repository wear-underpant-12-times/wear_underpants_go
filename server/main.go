package main

import (
	"errors"
	"flag"
	"io"
	"log"
	"net"

	"../utils"
)

var (
	port = "8082"
	h    = false
	// Commands = []string{"CONNECT", "BIND", "UDP ASSOCIATE"}
	AddrType = []string{"", "IPv4", "", "Domain", "IPv6"}

	errAddrType      = errors.New("socks addr type not supported")
	errVer           = errors.New("socks version not supported")
	errMethod        = errors.New("socks only support noauth method")
	errAuthExtraData = errors.New("socks authentication get extra data")
	errReqExtraData  = errors.New("socks request get extra data")
	errCmd           = errors.New("socks only support connect command")
)

func shake(conn net.Conn) (target string, err error) {
	lenBuf := make([]byte, 1)
	if _, err = io.ReadFull(conn, lenBuf); err != nil {
		log.Println(err)
		return "", err
	}
	buf := make([]byte, int(lenBuf[0]))
	if _, err = io.ReadFull(conn, buf); err != nil {
		log.Println("get full failed")
		return "", err
	}
	addr, err := utils.UnPackData(buf)
	if err != nil {
		return "", err
	}
	return string(addr), nil
}

func handConn(conn net.Conn) {
	defer func() {
		conn.Close()
		// log.Println("close client connection:", conn)
	}()
	addr, err := shake(conn)
	if err != nil {
		log.Println("shake error", err)
		return
	}
	log.Println(conn.RemoteAddr(), "->", addr)
	remoteConn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println("connect server error:", addr, err)
		return
	}
	defer func() {
		remoteConn.Close()
		// log.Println("close remote connection:", addr)
	}()
	go utils.NetDecodeCopy(conn, remoteConn)
	utils.NetEncodeCopy(remoteConn, conn)
}

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC | log.Lshortfile)
	flag.BoolVar(&h, "h", false, "this help")
	flag.StringVar(&port, "p", "8082", "port")
	flag.Parse()
}

func main() {
	if h {
		flag.Usage()
		return
	}
	l, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		panic(err)
	}
	log.Printf("start server on %s ...", port)
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handConn(conn)
	}
}
