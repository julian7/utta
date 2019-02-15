package connector

import "net"

// TCPConnector is a Connector implementation, providing pure
// TCP-based connection
type TCPConnector struct{}

// NewTCPConnector returns a new TCP connector
func NewTCPConnector() (Connector, error) {
	return &TCPConnector{}, nil
}

// Setup takes a net.Conn, and returns it verbatim.
func (*TCPConnector) Setup(conn net.Conn) (net.Conn, error) {
	return conn, nil
}
