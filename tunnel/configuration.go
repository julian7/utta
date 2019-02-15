package tunnel

import (
	"errors"
	"flag"
)

type configuration struct {
	listenAddr  string
	listenCert  string
	listenKey   string
	listenCA    string
	connectAddr string
	connectCert string
	connectKey  string
	serverName  string
	tls         bool
	proxy       string
	sshTunnel   string
	sshUser     string
	sshKey      string
}

func GetConfiguration() (*configuration, error) {
	conf := &configuration{}
	flag.StringVar(&conf.listenAddr, "listen", ":8080", "Listen port")
	flag.StringVar(&conf.listenCert, "lcert", "", "TLS certificate bundle for listening port (optional)")
	flag.StringVar(&conf.listenKey, "lkey", "", "TLS key for listening port (optional, default to -lcert)")
	flag.StringVar(&conf.listenCA, "lca", "", "mTLS accepted CA certs for listening port (turns on mTLS, optional)")
	flag.StringVar(&conf.connectAddr, "connect", "", "Connect port")
	flag.StringVar(&conf.connectCert, "ccert", "", "TLS certificate for connection (optional, sets -tls)")
	flag.StringVar(&conf.connectKey, "ckey", "", "TLS key for connection (optional, default to -ccert)")
	flag.StringVar(&conf.serverName, "servername", "", "Server name (only for TLS when connect name is different than SNI)")
	flag.BoolVar(&conf.tls, "tls", false, "TLS connection")
	flag.StringVar(&conf.proxy, "proxy", "", "Proxy host:port (default: no proxy)")
	flag.StringVar(&conf.sshTunnel, "sshtunnel", "", "SSH server host:port (default: no tunnel)")
	flag.StringVar(&conf.sshUser, "sshuser", "", "SSH username (required for SSH tunnel)")
	flag.StringVar(&conf.sshKey, "sshkey", "", "SSH private key file (required for SSH tunnel)")
	flag.Parse()

	if len(conf.connectAddr) < 1 {
		return nil, errors.New("providing -connect is required")
	}
	if len(conf.connectCert) >= 1 {
		conf.tls = true
	}
	if len(conf.sshTunnel) >= 1 {
		if conf.sshUser == "" {
			return nil, errors.New("please provide -sshuser")
		}
		if conf.sshKey == "" {
			return nil, errors.New("please provide -sshkey")
		}
	}

	return conf, nil
}
