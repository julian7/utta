package tunnel

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"net"
)

type Listener interface {
	Address() string
	Listen() (net.Listener, error)
}

type LocalListener struct {
	Addr     string
	CertFile string
	KeyFile  string
	CAFile   string
}

func NewListener(addr, cert, key, ca string) *LocalListener {
	return &LocalListener{
		Addr:     addr,
		CertFile: cert,
		KeyFile:  key,
		CAFile:   ca,
	}
}

func (l *LocalListener) Address() string {
	return l.Addr
}

func (l *LocalListener) Listen() (net.Listener, error) {
	ln, err := net.Listen("tcp", l.Addr)
	if err != nil {
		return nil, err
	}

	if l.CertFile == "" {
		return ln, nil
	}

	if l.KeyFile == "" {
		l.KeyFile = l.CertFile
	}
	certbundle, err := tls.LoadX509KeyPair(l.CertFile, l.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("cannot load cert keypair: %w", err)
	}

	config := &tls.Config{Certificates: []tls.Certificate{certbundle}}

	if l.CAFile != "" {
		caFile, err := ioutil.ReadFile(l.CAFile)
		if err != nil {
			return nil, fmt.Errorf("cannot read CA file: %w", err)
		}

		config.RootCAs = x509.NewCertPool()
		config.ClientCAs.AppendCertsFromPEM(caFile)
		config.ClientAuth = tls.RequireAndVerifyClientCert
	}
	config.BuildNameToCertificate()

	return tls.NewListener(ln, config), nil
}
