package connector

import (
	"io/ioutil"
	"net"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

type sshConnector struct {
	connect string
	tunnel  string
	client  *ssh.Client
	config  *ssh.ClientConfig
}

// NewSSHConnector sets up a new sshDialer for transmitting data through SSH tunnel
func NewSSHConnector(tunnel, user, key string) (Connector, error) {
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
	return &sshConnector{tunnel: tunnel, config: config}, nil
}

func (c *sshConnector) Setup(conn net.Conn) (net.Conn, error) {
	cconn, chans, reqs, err := ssh.NewClientConn(conn, c.connect, c.config)
	if err != nil {
		return nil, err
	}
	client := ssh.NewClient(cconn, chans, reqs)
	return client.Dial("tcp", c.tunnel)
}
