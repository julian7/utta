package remote

import (
	"time"

	"github.com/julian7/utta/tunnel"
	"go.uber.org/zap"
)

const (
	RETRY_SLEEP      = 5 * time.Second
	FLAPPING_TIMEOUT = 30 * time.Second
)

type Action struct {
	connection *tunnel.Connection
	sshlisten  string
	logger     *zap.Logger
	breaker    uint
	sleep      time.Duration
	counter    uint
}

func New(conn *tunnel.Connection, logger *zap.Logger, sshlisten string, breaker uint, sleep time.Duration) *Action {
	return &Action{
		connection: conn,
		sshlisten:  sshlisten,
		logger:     logger,
		breaker:    breaker,
		sleep:      sleep,
		counter:    0,
	}
}

func (a *Action) wait(timeToError time.Duration) {
	a.logger.Info("wait", zap.Duration("tte", timeToError))

	if timeToError < FLAPPING_TIMEOUT {
		a.counter += 1

		if a.counter%a.breaker == 0 {
			a.logger.Error(
				"circuit breaker",
				zap.Duration("sleep", a.sleep),
			)
			time.Sleep(a.sleep)
			a.counter = 0

			return
		}
	} else {
		a.counter = 0
	}

	time.Sleep(RETRY_SLEEP)
}

func (a *Action) Loop(tun *tunnel.Connection) {
	for {
		start := time.Now()
		conn, err := a.connection.Dial()
		if err != nil {
			a.logger.Error("connection error", zap.Error(err))
			a.wait(time.Now().Sub(start))
			continue
		}

		listener := tunnel.NewRemoteListener(
			a.sshlisten,
			a.connection,
		)

		start = time.Now()
		err = tunnel.NewTunnel(a.logger, listener).Run(tun)
		if err != nil {
			a.logger.Error("SSH tunnel error", zap.Error(err))
			_ = conn.Close()
			a.wait(time.Now().Sub(start))
			continue
		}
	}
}
