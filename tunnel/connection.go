package tunnel

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"github.com/julian7/utta/tunnel/connector"
	"github.com/julian7/utta/tunnel/dialer"
	"github.com/pkg/errors"
)

type connectionConfig struct {
	config  *Configuration
	dial    dialer.Dialer
	connect connector.Connector
	sshtun  *connector.SSHConnector
}

func (conf *connectionConfig) handleConn(downstream net.Conn) {
	defer downstream.Close()

	upstream, err := conf.dial.Dial()
	if err != nil {
		log.Printf("dial failed: %v", err)
		return
	}

	defer upstream.Close()
	log.Printf("Connected to %s", upstream.RemoteAddr())

	upstream, err = conf.connect.Setup(upstream)
	if err != nil {
		log.Printf("cannot set up TLS: %v", err)
		return
	}
	done := make(chan bool)
	if conf.sshtun != nil {
		upstream, err = conf.sshtun.Setup(upstream)
		if err != nil {
			log.Printf("cannot build SSH: %v", err)
			return
		}
		log.Printf("Built up SSH connection to %s through %s", conf.sshtun.Tunnel, conf.sshtun.Addr)
		go func(stream net.Conn, done <-chan bool) {
			select {
			case <-done:
				if stream != nil {
					stream.Close()
				}
				return
			}
		}(conf.sshtun.Connect, done)
	}
	go datapipe(downstream, upstream, "received", done)
	go datapipe(upstream, downstream, "transmitted", done)
	<-done
}

func datapipe(dst io.WriteCloser, src io.ReadCloser, direction string, done chan<- bool) {
	var errstr string
	b, err := io.Copy(io.Writer(dst), io.Reader(src))
	if err != nil {
		errOrig := errors.Cause(err)
		if strings.Contains(errOrig.Error(), "use of closed network connection") {
			err = nil
		}
	}
	if err != nil {
		errstr = fmt.Sprintf(" Error: %v vs %v", err, io.EOF)
	}
	log.Printf("%d bytes %s.%s", b, direction, errstr)
	done <- true
}
