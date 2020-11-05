package tunnel

import (
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/julian7/utta/tunnel/connector"
	"github.com/julian7/utta/tunnel/dialer"
	"github.com/julian7/utta/uuid"
	"github.com/pkg/errors"
)

type Connection struct {
	log.Logger
	dial    dialer.Dialer
	connect connector.Connector
	sshtun  *connector.SSHConnector
}

func NewConnection(logger log.Logger, addr, proxy string) *Connection {
	return &Connection{
		Logger: logger,
		dial:   dialer.NewDialer(5*time.Second, proxy, addr),
	}
}

func (cnx *Connection) SetupTCPConnector() error {
	conn, err := connector.NewTCPConnector()
	if err != nil {
		return err
	}

	cnx.connect = conn

	return nil
}

func (cnx *Connection) SetupTLSConnector(serverName, certFile, keyFile string) error {
	connect, err := connector.NewTLSConnector(serverName, certFile, keyFile)
	if err != nil {
		return err
	}

	cnx.connect = connect

	return nil
}

func (cnx *Connection) SetupSSHConnector(addr, tunnel, user, keyFile string) error {
	sshtun, err := connector.NewSSHConnector(addr, tunnel, user, keyFile)
	if err != nil {
		return err
	}

	cnx.sshtun = sshtun

	return nil
}

func (cnx *Connection) HasSSH() bool {
	return cnx.sshtun != nil
}

func (cnx *Connection) ListenSSH(addr string) (net.Listener, error) {
	return cnx.sshtun.Client.Listen("tcp", addr)
}

func (cnx *Connection) handleConn(downstream net.Conn) {
	defer downstream.Close()

	uuID, err := uuid.New()
	if err != nil {
		cnx.Logger.Log(
			"level", "error",
			"msg", "cannot generate UUID for connection",
			"connection", downstream.RemoteAddr().String(),
			"err", err,
		)
		return
	}

	connlog := log.With(cnx.Logger, "session", uuID.String())

	connlog.Log(
		"msg", "Handling connection",
		"remote", downstream.RemoteAddr().String(),
	)

	upstream, err := cnx.Dial()

	if err != nil {
		_ = connlog.Log("level", "error", "msg", "dial error", "err", err)

		return
	}

	defer upstream.Close()

	connlog.Log("msg", "connection established", "addr", upstream.RemoteAddr())

	done := make(chan bool)
	if cnx.sshtun != nil {
		connlog.Log(
			"msg", "SSH connection established",
			"addr", cnx.sshtun.Tunnel,
			"through", cnx.sshtun.Addr,
		)

		go func(stream net.Conn, done <-chan bool) {
			<-done
			if stream != nil {
				stream.Close()
				connlog.Log("msg", "SSH connection closed")
			}
		}(cnx.sshtun.Connect, done)
	}

	go datapipe(connlog, downstream, upstream, "received", done)
	go datapipe(connlog, upstream, downstream, "transmitted", done)
	<-done
}

func (cnx *Connection) Dial() (net.Conn, error) {
	upstream, err := cnx.dial.Dial()
	if err != nil {
		return nil, fmt.Errorf("dialing to remote: %w", err)
	}

	if cnx.connect == nil {
		return nil, errors.New("no connector provided for dial")
	}
	upstream, err = cnx.connect.Setup(upstream)
	if err != nil {
		return nil, fmt.Errorf("setting up TLS: %w", err)
	}
	if cnx.sshtun != nil {
		upstream, err = cnx.sshtun.Setup(upstream)
		if err != nil {
			return nil, fmt.Errorf("setting up SSH connection: %w", err)
		}
	}

	return upstream, nil
}

func datapipe(logger log.Logger, dst io.WriteCloser, src io.ReadCloser, direction string, done chan bool) {
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
	logger.Log(
		"msg", "connection summary",
		direction, b,
		"err", errstr,
	)
	done <- true
}
