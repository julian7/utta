package dialer

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// ProxyDialer is a Dialer implementation, to send packets through HTTP proxy
type ProxyDialer struct {
	timeout time.Duration
	proxy   string
	connect string
}

// NewProxyDialer sets up a new ProxyDialer
func NewProxyDialer(timeout time.Duration, proxy, connect string) *ProxyDialer {
	return &ProxyDialer{timeout: timeout, proxy: proxy, connect: connect}
}

// Dial connects to HTTP proxy, and requests remote connection using
// CONNECT method to the final destination. Returns net.Conn on success, error on error.
func (d *ProxyDialer) Dial() (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", d.proxy, d.timeout)
	if err != nil {
		return nil, errors.Wrap(err, "cannot dial proxy")
	}
	_, err = conn.Write([]byte(fmt.Sprintf("CONNECT %s HTTP/1.1\r\n\r\n", d.connect)))
	if err != nil {
		return nil, errors.Wrap(err, "cannot send CONNECT")
	}
	readbuf := make([]byte, 256)
	_, err = conn.Read(readbuf)
	if err != nil {
		return nil, errors.Wrap(err, "cannot read from proxy")
	}

	if !strings.HasPrefix(string(readbuf), "HTTP/1.1 200 ") {
		return nil, errors.Wrapf(err, "unknown proxy answer: %s", string(readbuf))
	}
	return conn, err
}
