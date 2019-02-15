package connector

import "net"

// Connector interface can set up a net.Conn for certain media
type Connector interface {
	Setup(net.Conn) (net.Conn, error)
}
