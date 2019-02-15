package tunnel

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/julian7/utta/tunnel/connector"
	"github.com/julian7/utta/tunnel/dialer"
	"github.com/pkg/errors"
)

// Tunnel sets up the tunnel based on configuration settings
func (config *configuration) Tunnel() error {
	conf, err := config.configureTunnel()
	if err != nil {
		log.Fatalf("cannot configure tunnel: %v", err)
	}

	ln, err := net.Listen("tcp", config.listenAddr)
	if err != nil {
		log.Fatalf("cannot listen on port: %v", err)
	}
	if len(config.listenCert) > 0 {
		ln, err = tlsListener(ln, config.listenCert, config.listenKey, config.listenCA)
	}

	defer ln.Close()

	rand.Seed(time.Now().Unix())

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Cannot listen: %v", err)
			continue
		}
		go conf.handleConn(conn)
	}
}

func tlsListener(l net.Listener, cert, key, ca string) (net.Listener, error) {
	if key == "" {
		key = cert
	}
	certbundle, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		return nil, errors.Wrap(err, "cannot load cert keypair")
	}
	config := &tls.Config{Certificates: []tls.Certificate{certbundle}}

	if ca != "" {
		caFile, err := ioutil.ReadFile(ca)
		if err != nil {
			return nil, errors.Wrap(err, "cannot read CA file")
		}
		config.RootCAs = x509.NewCertPool()
		config.ClientCAs.AppendCertsFromPEM(caFile)
		config.ClientAuth = tls.RequireAndVerifyClientCert
	}
	config.BuildNameToCertificate()
	return tls.NewListener(l, config), nil
}

func (config *configuration) configureTunnel() (*connectionConfig, error) {
	var err error
	conf := &connectionConfig{config: config}

	if config.tls {
		conf.connect, err = connector.NewTLSConnector(config.serverName, config.connectCert, config.connectKey)
	} else {
		conf.connect, err = connector.NewTCPConnector()
	}
	if err != nil {
		return nil, errors.Wrap(err, "cannot set up connector")
	}

	if len(config.sshTunnel) > 0 {
		conf.sshtun, err = connector.NewSSHConnector(config.sshTunnel, config.sshUser, config.sshKey)
	}

	conf.dial = dialer.NewDialer(5*time.Second, config.proxy, config.connectAddr)

	return conf, nil
}
