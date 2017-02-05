package common

import "net"

type Conn net.Conn
type Listener net.Listener

func DialTCP(addr string) (net.Conn, error) {
	return net.Dial("tcp", addr)
}

func ListenTCP(addr string) (Listener, error) {
	return net.Listen("tcp", addr)
}
