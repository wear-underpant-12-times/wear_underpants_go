package client

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/wear_underpants/utils"
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
	localPort      = "8081"
	wearServerAddr = "23.106.157.33:8082"
	udpListener    = &net.UDPConn{}
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

func parseAddr(conn net.Conn) (host string, addrType int, cmdType int, err error) {
	buf := make([]byte, 256)
	var n int
	if n, err = io.ReadAtLeast(conn, buf, 5); err != nil {
		log.Println(err)
		return
	}
	if buf[0] != 0x05 {
		err = errVer
		return
	}
	// if buf[1] != 0x01 {
	// 	err = errCmd
	// 	return
	// }
	// log.Println("method type", buf[1])
	cmdType = int(buf[1])
	reqLen := -1
	addrType = int(buf[3])
	// log.Println("addr type", buf[3])
	switch buf[3] {
	case 1: //ipv4
		reqLen = 4 + (3 + 1 + 2) // ip addr len plus (header + port len)
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
			log.Println(err)
			return
		}
	} else {
		log.Println("get over")
		err = errReqExtraData
		return
	}
	switch buf[3] {
	case 1: //ipv4
		host = strconv.Itoa(int(buf[4])) + "." + strconv.Itoa(int(buf[5])) + "." + strconv.Itoa(int(buf[6])) + "." + strconv.Itoa(int(buf[7]))
	case 3:
		host = string(buf[5 : 5+buf[4]])
	}
	port := binary.BigEndian.Uint16(buf[reqLen-2 : reqLen])
	host = net.JoinHostPort(host, strconv.Itoa(int(port)))
	return
}

func pipWhenClose(conn net.Conn, target string, addrType int) error {
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
	if addrType == 1 {
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
	} else if addrType == 3 {
		header := []byte{5, 0, 0, 3}
		domain := strings.Split(target, ":")[0]
		bDomain := []byte(domain)
		dLen := utils.Int8ToBytes(len(bDomain))
		port, _ := strconv.Atoi(strings.Split(target, ":")[1])
		bPort := utils.Int16ToBytes(port)
		res, _ := utils.MergeBytes([][]byte{
			header,
			dLen,
			bDomain,
			bPort,
		})
		conn.Write(res)
	}

	defer serverConn.Close()
	go utils.NetEncodeCopy(conn, serverConn)
	utils.NetDecodeCopy(serverConn, conn)
	return nil
}

// func checkSum(msg []byte) uint16 {
// 	sum := 0
// 	for n := 1; n < len(msg)-1; n += 2 {
// 		sum += int(msg[n])*256 + int(msg[n+1])
// 	}
// 	sum = (sum >> 16) + (sum & 0xffff)
// 	sum += (sum >> 16)
// 	var ans = uint16(^sum)
// 	return ans
// }

func handConn(conn net.Conn) {
	defer conn.Close()
	if err := shake(conn); err != nil {
		log.Println("socks handshake error", err)
		return
	}

	host, addrType, cmdType, err := parseAddr(conn)
	if err != nil {
		log.Println("socks addr parse error")
		return
	}
	log.Println(host)
	if cmdType == 1 {
		err = pipWhenClose(conn, host, addrType)
		if err != nil {
			log.Println("wear server data communication failed", host)
			return
		}
	} else if cmdType == 3 {
		// https://blog.csdn.net/whatday/article/details/40183555
		// UDP穿透应答
		iPort, err := strconv.Atoi(localPort)
		udpAddr := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: iPort}
		req := make([]byte, 256)
		req[0] = 0x05
		req[1] = 0x00
		req[2] = 0x00
		req[3] = 0x01 //注意：按照标准协议，返回的应该是对应的address_type，但是实际测试发现，当address_type=3，也就是说是域名类型时，会出现卡死情况，但是将address_type该为1，则不管是IP类型和域名类型都能正常运行
		ip := udpAddr.IP.To4()
		pindex := 4
		for _, b := range ip {
			req[pindex] = b
			pindex++
		}
		req[pindex] = byte((udpAddr.Port >> 8) & 0xff)
		req[pindex+1] = byte(udpAddr.Port & 0xff)
		conn.Write(req[0 : pindex+2])
		// 获取到应用的udp数据
		data := make([]byte, 1024)
		n, remoteAddr, _ := udpListener.ReadFromUDP(data)
		log.Println("udp(dns data):", remoteAddr, data[:n])
		// log.Println("socks5 udp header:", data[:10])
		// 将data发送到wear server
		uaddr, err := net.ResolveUDPAddr("udp", wearServerAddr)
		if err != nil {
			log.Println(err)
			return
		}
		udpConn, err := net.DialUDP("udp", nil, uaddr)
		if err != nil {
			log.Println(err)
		}
		defer udpConn.Close()
		_, err = udpConn.Write(data[:n])
		if err != nil {
			fmt.Println("failed:", err)
		}
		rdata := make([]byte, 1024)
		// 返回udp数据，写入到data
		n, err = udpConn.Read(rdata)
		if err != nil {
			log.Println("get udp data failed", err)
		}
		// 发送到udp(基本都是dns)服务器
		// srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
		// dstAddr := &net.UDPAddr{IP: net.IP(data[4:8]), Port: int(binary.BigEndian.Uint16(data[8:10]))}
		// dnsConn, err := net.DialUDP("udp", srcAddr, dstAddr)
		// if err != nil {
		// 	log.Println(err)
		// }
		// defer dnsConn.Close()
		// dnsConn.Write(data[10:n])
		// n, err = dnsConn.Read(data)
		// if err != nil {
		// 	log.Println("get udp data failed", err)
		// }

		// 再将数据打包返回
		h := []byte{0, 0, 0, 1}
		addr := data[4:8]
		// ip = remoteAddr.IP.To4()
		// pindex = 0
		// for _, b := range ip {
		// 	addr[pindex] = b
		// 	pindex++
		// }
		// log.Println(addr)
		// dstPort := utils.Int16ToBytes(remoteAddr.Port)
		dstPort := data[8:10]
		decoded1, err := base64.StdEncoding.DecodeString(string(rdata[:n]))
		if err != nil {
			log.Println("decode udp data error:", err)
			return
		}
		res, _ := utils.MergeBytes([][]byte{
			h,
			addr,
			dstPort,
			decoded1,
		})
		_, err = udpListener.WriteToUDP(res, remoteAddr)
		log.Println("goRetrun udp(dns) data:", err, decoded1)
	}
}

// func init() {
// 	log.SetFlags(log.Ldate | log.Ltime | log.LUTC | log.Lshortfile)
// 	flag.BoolVar(&h, "h", false, "this help")
// 	flag.StringVar(&localProxyPort, "p", "8088", "local socks5 port")
// 	flag.StringVar(&wearServerAddr, "addr", "127.0.0.1:8082", "remote addr:port")
// 	flag.Parse()
// }

// ss
func StartClient(localProxyPort string, remoteAddr string) {
	wearServerAddr = remoteAddr
	localPort = localProxyPort
	start()
}

func start() {
	iPort, err := strconv.Atoi(localPort)
	if err != nil {
		panic(err)
	}
	l, err := net.Listen("tcp", "127.0.0.1:"+localPort)
	if err != nil {
		panic(err)
	}
	udpListener, err = net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: iPort})
	if err != nil {
		panic(err)
	}
	log.Printf("start client on %s ...", localPort)
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

// func main() {
// 	if h {
// 		flag.Usage()
// 		return
// 	}
// 	start(localProxyPort)
// }
