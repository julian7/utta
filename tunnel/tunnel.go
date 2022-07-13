package tunnel

import (
	"io"
	"math/rand"
	"time"

	"go.uber.org/zap"
)

type Tunnel struct {
	*zap.Logger
	Listener
}

func NewTunnel(logger *zap.Logger, l Listener) *Tunnel {
	return &Tunnel{
		Logger:   logger,
		Listener: l,
	}
}

func (t *Tunnel) Run(cnx *Connection) error {
	ln, err := t.Listen()
	if err != nil {
		t.Error("cannot listen on port", zap.Error(err))
		return err
	}
	t.Info("Tunnel listening", zap.String("addr", t.Listener.Address()))

	defer ln.Close()

	rand.Seed(time.Now().Unix())

	for {
		conn, err := ln.Accept()
		if err != nil {
			if err == io.EOF {
				t.Warn("connection closed")
				break
			}
			t.Error("error in listen", zap.Error(err))
			continue
		}
		go cnx.handleConn(conn)
	}

	return nil
}
