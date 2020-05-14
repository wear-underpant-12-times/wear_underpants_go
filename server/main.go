package main

import (
	"errors"
	"flag"
	"io"
	"log"
	"net"
	"time"

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
	buf := make([]byte, 258)
	var n int
	if n, err = io.ReadAtLeast(conn, buf, 1); err != nil {
		log.Println(err)
		return
	}

	dmLen := int(buf[0])
	msgLen := dmLen + 1
	if n == msgLen {
	} else if n < msgLen {
		if _, err = io.ReadFull(conn, buf[n:msgLen]); err != nil {
			log.Println("get full failed")
			return "", err
		}
	} else {
		log.Printf("dmLen %v, getLen %v\n", dmLen, n)
		// return errors.New("auth error")
	}
	addr, err := utils.UnPackData(buf[1 : 1+dmLen])
	if err != nil {
		return "", err
	}
	return string(addr), nil
}

func handConn(conn net.Conn) {
	defer conn.Close()
	addr, err := shake(conn)
	if err != nil {
		log.Println("shake error", err)
		return
	}
	log.Println(addr)
	remoteConn, err := net.DialTimeout("tcp", addr, time.Duration(time.Second*15))
	if err != nil {
		log.Println("connect server error:", addr, err)
		return
	}
	defer remoteConn.Close()
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
