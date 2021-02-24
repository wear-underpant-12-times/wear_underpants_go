package mobile

import (
	"strconv"
	"time"

	"github.com/eycorsican/go-tun2socks/core"
	"github.com/eycorsican/go-tun2socks/proxy/socks"
	"github.com/wear_underpants/client"
)

func StartServer(localPort string, addr string) {
	client.StartClient(localPort, addr)
}

type PacketFlow interface {
	WritePacket(packet []byte)
}

var lwipStack core.LWIPStack

func InputPacket(data []byte) {
	lwipStack.Write(data)
}

func Stop() {
	client.Stop()
	lwipStack.Close()
}

func StartSocks(packetFlow PacketFlow, proxyHost string, proxyPort int, addr string) {
	if packetFlow != nil {
		lwipStack = core.NewLWIPStack()
		core.RegisterTCPConnHandler(socks.NewTCPHandler(proxyHost, uint16(proxyPort)))
		core.RegisterUDPConnHandler(socks.NewUDPHandler(proxyHost, uint16(proxyPort), 2*time.Minute))
		core.RegisterOutputFn(func(data []byte) (int, error) {
			packetFlow.WritePacket(data)
			return len(data), nil
		})
		go client.StartClient(strconv.Itoa(proxyPort), addr)
	}
}
