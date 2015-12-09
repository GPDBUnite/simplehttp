package unite

import (
	"net"
	"fmt"
)

func UDPLogServer(port int, host string) {

	addr := net.UDPAddr{
		Port: port,
		IP: net.ParseIP(host),
    }
	//fmt.Println(addr)
    conn, err := net.ListenUDP("udp", &addr)

    if err != nil {
        fmt.Println(err)
		return
    }
    defer conn.Close()
	for {
		buf := make([]byte, 1024)
		rlen, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Print(string(buf[0:rlen]))
		}
	}

}
