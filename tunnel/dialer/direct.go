package dialer

import (
	"net"
	"time"
)

// DirectDialer is an implementation of Dialer, to dial directly
type DirectDialer struct {
	timeout time.Duration
	connect string
}

// NewDirectDialer is creating a new DirectDialer
func NewDirectDialer(timeout time.Duration, connect string) *DirectDialer {
	return &DirectDialer{timeout: timeout, connect: connect}
}

// Dial calls net.DialTimeout with provided data
func (d *DirectDialer) Dial() (net.Conn, error) {
	return net.DialTimeout("tcp", d.connect, d.timeout)
}
