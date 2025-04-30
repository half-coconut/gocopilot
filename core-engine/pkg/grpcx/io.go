package grpcx

import "net"

func GetOutBoundIP() string {
	// DNS 的地址，国内可以用 114.114.114.114
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return ""
	}
	defer func(conn net.Conn) {
		err = conn.Close()
		if err != nil {
			panic(err)
		}
	}(conn)

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}
