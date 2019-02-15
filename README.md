# Universal Travel TCP Adapter (UTTA)

You want a TCP tunnel, no matter what. You know your packet will go through dangerous places. Therefore, the best bet is to take the Universal Travel TCP Adapter to make sure your packets will have a safe journey to destination.

## What does it do?

UTTA is capable of

* listening to a TCP port, potentially requiring (m)TLS,
* opening a TCP connection to the destination when it receives a connection,
* sending this connection through HTTP proxy,
* set up a tunnel through an SSH connection,
* set up TLS for the outgoing connection.

## Usage

Currently, the app has a single CLI, with a couple of options:

```text
  -ccert string
    	TLS certificate for connection (optional, sets -tls)
  -ckey string
    	TLS key for connection (optional, default to -ccert)
  -connect string
    	Connect port
  -lca string
    	mTLS accepted CA certs for listening port (turns on mTLS, optional)
  -lcert string
    	TLS certificate bundle for listening port (optional)
  -listen string
    	Listen port (default ":8080")
  -lkey string
    	TLS key for listening port (optional, default to -lcert)
  -proxy string
    	Proxy host:port (default: no proxy)
  -servername string
    	Server name (only for TLS when connect name is different than SNI)
  -sshkey string
    	SSH private key file (required for SSH tunnel)
  -sshtunnel string
    	SSH server host:port (default: no tunnel)
  -sshuser string
    	SSH username (required for SSH tunnel)
  -tls
    	TLS connection
```

## Any issues?

Open a ticket, perhaps a pull request. We support [GitHub Flow](https://guides.github.com/introduction/flow/)
