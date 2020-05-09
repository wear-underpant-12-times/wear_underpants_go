package main

import (
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"../utils"
)

var (
	// Commands = []string{"CONNECT", "BIND", "UDP ASSOCIATE"}
	AddrType = []string{"", "IPv4", "", "Domain", "IPv6"}

	errAddrType      = errors.New("socks addr type not supported")
	errVer           = errors.New("socks version not supported")
	errMethod        = errors.New("socks only support noauth method")
	errAuthExtraData = errors.New("socks authentication get extra data")
	errReqExtraData  = errors.New("socks request get extra data")
	errCmd           = errors.New("socks only support connect command")

	h              = false
	localProxyPort = "8081"
	wearServerAddr = "23.106.157.33:8082"
)

func shake(conn net.Conn) (err error) {
	buf := make([]byte, 258)
	var n int
	if n, err = io.ReadAtLeast(conn, buf, 2); err != nil {
		return
	}

	if buf[0] != 0x05 {
		return errors.New("socks version not support")
	}

	nmethod := int(buf[1])
	msgLen := nmethod + 2
	if n == msgLen {

	} else if n < msgLen {
		if _, err = io.ReadFull(conn, buf[n:msgLen]); err != nil {
			return
		}
	} else {
		return errors.New("auth error")
	}
	_, err = conn.Write([]byte{0x05, 0})
	return
}

func parseAddr(conn net.Conn) (host string, err error) {
	buf := make([]byte, 256)
	var n int
	if n, err = io.ReadAtLeast(conn, buf, 5); err != nil {
		return
	}
	if buf[0] != 0x05 {
		err = errVer
		return
	}
	if buf[1] != 0x01 {
		err = errCmd
		return
	}

	reqLen := -1
	switch buf[3] {
	case 1: //ipv4
		fmt.Println("not support ipv4")
		return
	case 3: //domain
		reqLen = int(buf[4]) + (3 + 1 + 1 + 2) //(3 + 1 + 1 + 2) 3 + 1addrType + 1addrLen + 2port, plus addrLen
	case 4: //ipv6
		fmt.Println("not support ipv6")
		return
	default:
		err = errAddrType
		return
	}

	if n == reqLen {

	} else if n < reqLen {
		if _, err = io.ReadFull(conn, buf[n:reqLen]); err != nil {
			return
		}
	} else {
		err = errReqExtraData
		return
	}

	switch buf[3] {
	case 3:
		host = string(buf[5 : 5+buf[4]])
	}
	port := binary.BigEndian.Uint16(buf[reqLen-2 : reqLen])
	host = net.JoinHostPort(host, strconv.Itoa(int(port)))
	return
}

func pipWhenClose(conn net.Conn, target string) error {
	// 与wear服务器建立连接
	serverConn, err := net.DialTimeout("tcp", wearServerAddr, time.Duration(time.Second*15))
	if err != nil {
		log.Println("connect wear server error:", target, err)
		return err
	}
	bAddr := []byte(target)
	packbAddr, err := utils.PackHeader(bAddr)
	if err != nil {
		return err
	}
	serverConn.Write(packbAddr)
	tcpAddr := serverConn.LocalAddr().(*net.TCPAddr)
	req := make([]byte, 256)
	req[0] = 0x05
	req[1] = 0x00
	req[2] = 0x00
	req[3] = 0x01 //注意：按照标准协议，返回的应该是对应的address_type，但是实际测试发现，当address_type=3，也就是说是域名类型时，会出现卡死情况，但是将address_type该为1，则不管是IP类型和域名类型都能正常运行
	ip := tcpAddr.IP.To4()
	pindex := 4
	for _, b := range ip {
		req[pindex] = b
		pindex++
	}
	req[pindex] = byte((tcpAddr.Port >> 8) & 0xff)
	req[pindex+1] = byte(tcpAddr.Port & 0xff)
	conn.Write(req[0 : pindex+2])
	defer serverConn.Close()
	go utils.NetEncodeCopy(conn, serverConn)
	utils.NetDecodeCopy(serverConn, conn)
	return nil
}

func handConn(conn net.Conn) {
	defer conn.Close()
	if err := shake(conn); err != nil {
		log.Println("socks handshake error")
		return
	}

	host, err := parseAddr(conn)
	if err != nil {
		log.Println("socks addr parse error")
		return
	}
	log.Println(host)
	err = pipWhenClose(conn, host)
	if err != nil {
		log.Println("wear server data communication failed", host)
		return
	}
}

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC | log.Lshortfile)
	flag.BoolVar(&h, "h", false, "this help")
	flag.StringVar(&localProxyPort, "p", "8081", "local socks5 port")
	flag.StringVar(&wearServerAddr, "addr", "127.0.0.1:8082", "remote addr:port")
	flag.Parse()

}

func main() {
	if h {
		flag.Usage()
		return
	}
	l, err := net.Listen("tcp", "127.0.0.1:"+localProxyPort)
	if err != nil {
		panic(err)
	}
	log.Printf("start client on %s ...", localProxyPort)
	log.Printf("wear server %s", wearServerAddr)
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handConn(conn)
	}
}
