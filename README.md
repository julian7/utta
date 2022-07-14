# Universal Travel TCP Adapter (UTTA)

You want a TCP tunnel, no matter what. You know your packet will go through dangerous places. Therefore, the best bet is to take the Universal Travel TCP Adapter to make sure your packets will have a safe journey to destination.

## What does it do?

UTTA is capable of

* listening to a TCP port, potentially requiring (m)TLS,
* opening a TCP connection to the destination when it receives a connection,
* sending this connection through HTTP proxy,
* set up a tunnel through an SSH connection,
* set up TLS for the outgoing connection,
* or, building up a TCP connection to a SSH endpoint, and provide remote proxy

## Usage

The application has two modes of running: local and remote. Local provides a
local listening port, which connects to remote TCP / TLS / SSH service.
Remote keeps a SSH connection up, listens at the remote server, forwarding all
connections to a local service.

```text
NAME:
   utta - Universal Travel TCP Adapter

USAGE:
   utta [global options] command [command options] [arguments...]

COMMANDS:
   local    create locally listening tunnel
   remote   create remotely listening tunnel
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h                   show help (default: false)
   --log-level value, -l value  Log level (default: info; values: debug, info, warn, error, panic, fatal) [$UTTA_LOG_LEVEL]
   --version, -v                print the version (default: false)
```

Options for local operations:

```text
NAME:
   utta local - create locally listening tunnel

USAGE:
   utta local [command options] [arguments...]

OPTIONS:
   --ccert value       Client TLS cert for connect [$UTTA_CONNECT_CERT]
   --ckey value        Client TLS private key for connect [$UTTA_CONNECT_KEY]
   --connect value     Connect port [$UTTA_CONNECT]
   --lca value         Server TLS CA cert bundle [$UTTA_LISTEN_CA]
   --lcert value       Server TLS cert for listen [$UTTA_LISTEN_CERT]
   --listen value      Listen port (default: ":8080") [$UTTA_LISTEN]
   --lkey value        Server TLS private key for listen [$UTTA_LISTEN_KEY]
   --proxy value       HTTP proxy host:port (default: no proxy) [$UTTA_PROXY]
   --servername value  Server name for TLS connect with SNI [$UTTA_CONNECT_SERVERNAME]
   --sshkey value      SSH key for tunnel [$UTTA_SSH_KEY]
   --sshtunnel value   SSH server host:port (default: no tunnel through SSH) [$UTTA_SSH_TUNNEL]
   --sshuser value     SSH username for tunnel [$UTTA_SSH_USER]
   --tls               Connect with TLS (default: false) [$UTTA_CONNECT_TLS]
```

In this mode, UTTA listens on a local port (TLS/mTLS is optional), which builds up a connection on demand. It connects to a remote port (TLS is optional), traversing a HTTP proxy if needed (no proxy authentication implemented). Then, if sshtunnel is provided, it treats remote connect port as an SSH server, and connects to it with provided SSH user and key. Lastly, it establishes a forwarding connection on top of SSH.

Options for remote operations:

```text
NAME:
   utta remote - create remotely listening tunnel

USAGE:
   utta remote [command options] [arguments...]

OPTIONS:
   --breaker value     Circuit breaker: taking a break after # attempts (default: 3) [$UTTA_BREAKER]
   --ccert value       Client TLS cert for connect [$UTTA_CONNECT_CERT]
   --ckey value        Client TLS private key for connect [$UTTA_CONNECT_KEY]
   --connect value     Connect port [$UTTA_CONNECT]
   --proxy value       HTTP proxy host:port (default: no proxy) [$UTTA_PROXY]
   --servername value  Server name for TLS connect with SNI [$UTTA_CONNECT_SERVERNAME]
   --sleep value       Sleep between circuit breaks (default: 30m0s) [$UTTA_SLEEP]
   --sshconnect value  SSH local target port [$UTTA_SSH_CONNECT]
   --sshkey value      SSH key for tunnel [$UTTA_SSH_KEY]
   --sshlisten value   SSH remote listening port [$UTTA_SSH_LISTEN]
   --sshuser value     SSH username for tunnel [$UTTA_SSH_USER]
   --tls               Connect with TLS (default: false) [$UTTA_CONNECT_TLS]
```

In this mode, UTTA establishes (and restarts, if needed) a connection to a remote port (TLS is optional), traversing a HTTP proxy if needed (as with local mode, proxy authentication is not implemented). Then, it establishes an SSH connection with provided SSH user and key. Lastly, it establishes a remote port forwarding, listening at SSH endpoint, forwarding all connections to sshconnect host/port.

This mode has a circuit breaker. By default, it sleeps 5 seconds between connections. However, if a TCP (or the internal SSH) connection returns within 30 seconds three times in a row (configurable with `--breaker`), it will take a longer sleep (30 minutes by default, configurable with `--sleep`).

## Any issues?

Open a ticket, perhaps a pull request. We support [GitHub Flow](https://guides.github.com/introduction/flow/)
