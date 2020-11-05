package connector

import (
	"io/ioutil"
	"net"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

// SSHConnector is a Connector, which can connect using SSH through the original Connection
type SSHConnector struct {
	// Connect is the original connection
	Connect net.Conn
	// Addr is hostname:port of the SSH endpoint
	Addr string
	// Tunnel is hostname:port to be accessed through SSH
	Tunnel string
	Client *ssh.Client
	config *ssh.ClientConfig
}

// NewSSHConnector sets up a new sshDialer for transmitting data through SSH tunnel
func NewSSHConnector(addr, tunnel, user, key string) (*SSHConnector, error) {
	privkey, err := ioutil.ReadFile(key)
	if err != nil {
		return nil, errors.Wrap(err, "cannot read SSH private key")
	}
	signer, err := ssh.ParsePrivateKey(privkey)
	if err != nil {
		return nil, errors.Wrap(err, "invalid SSH private key format")
	}
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	return &SSHConnector{Addr: addr, Tunnel: tunnel, config: config}, nil
}

// Setup sets up the original net.Conn to connect through SSH
func (c *SSHConnector) Setup(conn net.Conn) (net.Conn, error) {
	cconn, chans, reqs, err := ssh.NewClientConn(conn, c.Addr, c.config)
	if err != nil {
		return nil, err
	}
	c.Client = ssh.NewClient(cconn, chans, reqs)

	if c.Tunnel != "" {
		return c.Client.Dial("tcp", c.Tunnel)
	}

	return conn, nil
}
