package server

import (
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
	"strconv"

	"github.com/wear_underpants/utils"
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

	udpListener = &net.UDPConn{}
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

func handleUDPConn(data []byte, n int, remoteAddr *net.UDPAddr) {
	log.Println("udp:", remoteAddr, "->", net.IP(data[4:8]))
	// log.Println("recieve udp request:", remoteAddr, data[:n])
	// 发送到udp(基本都是dns)服务器
	srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	dstAddr := &net.UDPAddr{IP: net.IP(data[4:8]), Port: int(binary.BigEndian.Uint16(data[8:10]))}
	dnsConn, err := net.DialUDP("udp", srcAddr, dstAddr)
	if err != nil {
		log.Println(err)
		return
	}
	defer dnsConn.Close()
	dnsConn.Write(data[10:n])
	// log.Println("sendUdp", data[10:n])
	// 返回
	n, err = dnsConn.Read(data)
	if err != nil {
		log.Println("get udp data failed", err)
	}
	// log.Println("getUdp", data[:n])
	// res := base64.StdEncoding.EncodeToString(data[:n])
	if err != nil {
		log.Println("pack udp data error:", data[:n])
	}

	// res, _ := utils.MergeBytes([][]byte{
	// 	[]byte{112},
	// 	data[:n],
	// })
	_, err = udpListener.WriteToUDP(data[:n], remoteAddr)
	// log.Println("return udp", err, data[:n])
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

// func init() {
// 	log.SetFlags(log.Ldate | log.Ltime | log.LUTC | log.Lshortfile)
// 	flag.BoolVar(&h, "h", false, "this help")
// 	flag.StringVar(&port, "p", "8082", "port")
// 	flag.Parse()
// }

// Start Server
func StartServer(port string) {
	iPort, err := strconv.Atoi(port)
	if err != nil {
		panic(err)
	}
	// tcp服务
	l, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		panic(err)
	}
	// udp服务，与tcp在同一端口
	udpListener, err = net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: iPort})
	if err != nil {
		panic(err)
	}
	log.Printf("start server on %s ...", port)
	data := make([]byte, 65535)
	go func() {
		for {
			n, remoteAddr, err := udpListener.ReadFromUDP(data)
			if err != nil {
				log.Println("failed to read UDP msg because of ", err.Error())
				continue
			}
			if n >= 10 && n <= 65535 {
				go handleUDPConn(data, n, remoteAddr)
			} else {
				log.Println("error: udp data length is less then 10.", n)
			}
		}
	}()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handConn(conn)
	}
}

// func main() {
// 	if h {
// 		flag.Usage()
// 		return
// 	}
// 	iPort, err := strconv.Atoi(port)
// 	if err != nil {
// 		panic(err)
// 	}
// 	// tcp服务
// 	l, err := net.Listen("tcp", "0.0.0.0:"+port)
// 	if err != nil {
// 		panic(err)
// 	}
// 	// udp服务，与tcp在同一端口
// 	udpListener, err = net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: iPort})
// 	if err != nil {
// 		panic(err)
// 	}
// 	log.Printf("start server on %s ...", port)
// 	go func() {
// 		for {
// 			handleUDPConn(udpListener)
// 		}
// 	}()
// 	for {
// 		conn, err := l.Accept()
// 		if err != nil {
// 			log.Println(err)
// 			continue
// 		}
// 		go handConn(conn)
// 	}
// }
