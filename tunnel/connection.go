package tunnel

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/julian7/utta/tunnel/connector"
	"github.com/julian7/utta/tunnel/dialer"
	"github.com/julian7/utta/uuid"
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

	uuID, err := uuid.New()
	sessionID := fmt.Sprintf("[%s] ", uuID.String())
	if err != nil {
		log.Printf(
			"Error: cannot generate UUID for connection %s: %v",
			downstream.RemoteAddr().String(),
			err,
		)
		return
	}
	connlog := log.New(os.Stdout, sessionID, log.LstdFlags)

	connlog.Printf(
		"Handling connection from %s", downstream.RemoteAddr().String(),
	)

	upstream, err := conf.dial.Dial()
	if err != nil {
		connlog.Printf("dial failed: %v", err)
		return
	}

	defer upstream.Close()
	connlog.Printf("Connected to %s", upstream.RemoteAddr())

	upstream, err = conf.connect.Setup(upstream)
	if err != nil {
		connlog.Printf("cannot set up TLS: %v", err)
		return
	}
	done := make(chan bool)
	if conf.sshtun != nil {
		upstream, err = conf.sshtun.Setup(upstream)
		if err != nil {
			connlog.Printf("cannot build SSH: %v", err)
			return
		}
		connlog.Printf(
			"Built up SSH connection to %s through %s",
			conf.sshtun.Tunnel,
			conf.sshtun.Addr,
		)
		go func(stream net.Conn, done <-chan bool) {
			select {
			case <-done:
				if stream != nil {
					stream.Close()
					connlog.Println("SSH connection closed")
				}
				return
			}
		}(conf.sshtun.Connect, done)
	}
	go datapipe(connlog, downstream, upstream, "received", done)
	go datapipe(connlog, upstream, downstream, "transmitted", done)
	<-done
}

func datapipe(logger *log.Logger, dst io.WriteCloser, src io.ReadCloser, direction string, done chan bool) {
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
	logger.Printf("%d bytes %s.%s", b, direction, errstr)
	done <- true
}
