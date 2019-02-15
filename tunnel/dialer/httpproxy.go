package dialer

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type proxyDialer struct {
	timeout time.Duration
	proxy   string
	connect string
}

func NewProxyDialer(timeout time.Duration, proxy, connect string) *proxyDialer {
	dialer := &proxyDialer{timeout: timeout, proxy: proxy, connect: connect}
	return dialer
}

func (d *proxyDialer) Dial() (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", d.proxy, d.timeout)
	if err != nil {
		return nil, errors.Wrap(err, "cannot dial proxy")
	}
	_, err = conn.Write([]byte(fmt.Sprintf("CONNECT %s HTTP/1.1\r\n\r\n", d.connect)))
	if err != nil {
		return nil, errors.Wrap(err, "cannot send CONNECT")
	}
	readbuf := make([]byte, 256)
	n, err := conn.Read(readbuf)
	if err != nil {
		return nil, errors.Wrap(err, "cannot read from proxy")
	}

	if !strings.HasPrefix(string(readbuf), "HTTP/1.1 200 ") {
		return nil, errors.Wrapf(err, "unknown proxy answer: %s", string(readbuf))
	}
	log.Printf("Read %d bytes: %s", n, string(readbuf))
	return conn, err
}
