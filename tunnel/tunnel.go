package tunnel

import (
	"io"
	"math/rand"
	"time"

	"github.com/go-kit/log"
)

type Tunnel struct {
	log.Logger
	Listener
}

func NewTunnel(logger log.Logger, l Listener) *Tunnel {
	return &Tunnel{
		Logger:   logger,
		Listener: l,
	}
}

func (t *Tunnel) Run(cnx *Connection) error {
	ln, err := t.Listen()
	if err != nil {
		_ = t.Log("level", "fatal", "msg", "cannot listen on port", "err", err)
		return err
	}
	_ = t.Log("msg", "Tunnel listening", "addr", t.Listener.Address())

	defer ln.Close()

	rand.Seed(time.Now().Unix())

	for {
		conn, err := ln.Accept()
		if err != nil {
			if err == io.EOF {
				_ = t.Log("level", "warn", "msg", "connection closed")
				break
			}
			_ = t.Log("level", "error", "msg", "error in listen", "err", err)
			continue
		}
		go cnx.handleConn(conn)
	}

	return nil
}
