//go:generate go run ../../internal/run/gen_versiontxt.go

package main

import (
	_ "embed"
	"fmt"

	"github.com/julian7/utta/tunnel"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//go:embed version.txt
var version string

type App struct {
	logger *zap.Logger
	level  *zap.AtomicLevel
}

func NewApp(l *zap.Logger, level *zap.AtomicLevel) *App {
	return &App{logger: l, level: level}
}

func (a *App) CommonFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "connect",
			Value:   "",
			Usage:   "Connect port",
			EnvVars: []string{"UTTA_CONNECT"},
		},
		&cli.StringFlag{
			Name:      "ccert",
			Value:     "",
			TakesFile: true,
			Usage:     "Client TLS cert for connect",
			EnvVars:   []string{"UTTA_CONNECT_CERT"},
		},
		&cli.StringFlag{
			Name:      "ckey",
			Value:     "",
			TakesFile: true,
			Usage:     "Client TLS private key for connect",
			EnvVars:   []string{"UTTA_CONNECT_KEY"},
		},
		&cli.StringFlag{
			Name:    "servername",
			Value:   "",
			Usage:   "Server name for TLS connect with SNI",
			EnvVars: []string{"UTTA_CONNECT_SERVERNAME"},
		},
		&cli.BoolFlag{
			Name:    "tls",
			Value:   false,
			Usage:   "Connect with TLS",
			EnvVars: []string{"UTTA_CONNECT_TLS"},
		},
		&cli.StringFlag{
			Name:    "proxy",
			Value:   "",
			Usage:   "HTTP proxy host:port (default: no proxy)",
			EnvVars: []string{"UTTA_PROXY"},
		},
		&cli.StringFlag{
			Name:    "sshuser",
			Value:   "",
			Usage:   "SSH username for tunnel",
			EnvVars: []string{"UTTA_SSH_USER"},
		},
		&cli.StringFlag{
			Name:      "sshkey",
			Value:     "",
			TakesFile: true,
			Usage:     "SSH key for tunnel",
			EnvVars:   []string{"UTTA_SSH_KEY"},
		},
	}
}

func (a *App) Command() *cli.App {
	return &cli.App{
		Name:    "utta",
		Usage:   "Universal Travel TCP Adapter",
		Version: version,
		Before:  a.setLogLevel,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "log-level",
				Aliases: []string{"l"},
				Value:   "",
				Usage:   "Log level (default: info; values: debug, info, warn, error, panic, fatal)",
				EnvVars: []string{"UTTA_LOG_LEVEL"},
			},
		},
		Commands: []*cli.Command{
			a.localCommand(),
			a.remoteCommand(),
		},
	}
}

func (a *App) setLogLevel(c *cli.Context) error {
	loglevel := c.String("log-level")
	if len(loglevel) > 0 {
		lvl, err := zapcore.ParseLevel(loglevel)
		if err != nil {
			return fmt.Errorf("parsing log-level: %w", err)
		}

		a.level.SetLevel(lvl)
		a.logger.Debug("set log level", zap.String("loglevel", loglevel))
	}

	return nil
}

func (a *App) GenericConnection(c *cli.Context) (*tunnel.Connection, error) {
	connection := tunnel.NewConnection(
		a.logger,
		c.String("connect"),
		c.String("proxy"),
	)
	if c.Bool("tls") || c.String("ccert") != "" {
		if err := connection.SetupTLSConnector(
			c.String("servername"),
			c.String("ccert"),
			c.String("ckey"),
		); err != nil {
			return nil, err
		}
	} else {
		if err := connection.SetupTCPConnector(); err != nil {
			return nil, err
		}
	}

	if c.String("sshuser") != "" {
		if err := connection.SetupSSHConnector(
			c.String("connect"),
			c.String("sshtunnel"),
			c.String("sshuser"),
			c.String("sshkey"),
		); err != nil {
			return nil, err
		}
	}

	return connection, nil
}
