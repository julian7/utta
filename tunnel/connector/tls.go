package connector

import (
	"crypto/tls"
	"net"

	"github.com/pkg/errors"
)

// TLSConnector is a Connector implementation to tunnel connections
// through TLS
type TLSConnector struct {
	servername string
	cert       tls.Certificate
}

// NewTLSConnector returns a Connector, which can be used for TLS connections.
func NewTLSConnector(servername, cert, key string) (Connector, error) {
	conn := &TLSConnector{servername: servername}
	if cert != "" {
		certbundle, err := tls.LoadX509KeyPair(cert, key)
		if err != nil {
			return nil, errors.Wrap(err, "cannot load client keypair")
		}
		conn.cert = certbundle
	}
	return conn, nil
}

// Setup takes a net.Conn, and implements a simple, non-checking TLS
// client with SNI requirements.
func (c *TLSConnector) Setup(conn net.Conn) (net.Conn, error) {
	config := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         c.servername,
	}
	return tls.Client(conn, config), nil
}
