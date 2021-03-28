package tunnel

import "net"

type RemoteListener struct {
	ListenOn  string
	ConnectTo *Connection
}

func NewRemoteListener(listen string, connect *Connection) *RemoteListener {
	return &RemoteListener{
		ListenOn:  listen,
		ConnectTo: connect,
	}
}

func (l *RemoteListener) Address() string {
	return l.ListenOn
}

func (l *RemoteListener) Listen() (net.Listener, error) {
	return l.ConnectTo.ListenSSH(l.ListenOn)
}
