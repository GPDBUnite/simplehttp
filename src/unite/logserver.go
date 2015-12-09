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
    conn, err := net.ListenUDP("udp", &addr)
    defer conn.Close()
    if err != nil {
        panic(err)
    }

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
