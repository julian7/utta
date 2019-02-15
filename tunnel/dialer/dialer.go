package dialer

import (
	"net"
	"time"
)

// Dialer is anything which can Dial
type Dialer interface {
	Dial() (net.Conn, error)
}

// NewDialer returns an appropriate dialer for the job
func NewDialer(t time.Duration, proxy string, connect string) Dialer {
	if len(proxy) > 0 {
		return NewProxyDialer(5*time.Second, proxy, connect)
	}
	return NewDirectDialer(5*time.Second, connect)
}
