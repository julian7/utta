package main

import (
	"errors"

	"github.com/julian7/utta/tunnel"
	"github.com/urfave/cli/v2"
)

func (a *App) localCommand() *cli.Command {
	return &cli.Command{
		Name:  "local",
		Usage: "create locally listening tunnel",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "listen",
				Value:    ":8080",
				Usage:    "Listen port",
				Required: true,
				EnvVars:  []string{"UTTA_LISTEN"},
			},
			&cli.StringFlag{
				Name:      "lcert",
				Value:     "",
				TakesFile: true,
				Usage:     "Server TLS cert for listen",
				EnvVars:   []string{"UTTA_LISTEN_CERT"},
			},
			&cli.StringFlag{
				Name:      "lca",
				Value:     "",
				TakesFile: true,
				Usage:     "Server TLS CA cert bundle",
				EnvVars:   []string{"UTTA_LISTEN_CA"},
			},
			&cli.StringFlag{
				Name:      "lkey",
				Value:     "",
				TakesFile: true,
				Usage:     "Server TLS private key for listen",
				EnvVars:   []string{"UTTA_LISTEN_KEY"},
			},
			&cli.StringFlag{
				Name:    "sshtunnel",
				Value:   "",
				Usage:   "SSH server host:port (default: no tunnel through SSH)",
				EnvVars: []string{"UTTA_SSH_TUNNEL"},
			},
		},
		Action: a.localAction,
		Before: a.localBefore,
	}
}

func (a *App) localAction(c *cli.Context) error {
	connection, err := a.GenericConnection(c)
	if err != nil {
		return err
	}

	listener := tunnel.NewListener(
		c.String("listen"),
		c.String("lcert"),
		c.String("lkey"),
		c.String("lca"),
	)

	return tunnel.NewTunnel(a.logger, listener).Run(connection)
}

func (a *App) localBefore(c *cli.Context) error {
	if c.String("connect") == "" {
		return errors.New("--connect is required")
	}
	if c.String("sshtunnel") != "" {
		if c.String("sshuser") == "" {
			return errors.New("--sshuser not provided")
		}
		if c.String("sshkey") == "" {
			return errors.New("--sshkey not provided")
		}
	}
	return nil
}
