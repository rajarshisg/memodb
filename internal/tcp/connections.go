package tcp

import "net"

type TCPConnection struct { 
	Connection net.Conn
	Initialized bool
}
var Connections []TCPConnection;