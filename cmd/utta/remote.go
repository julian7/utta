package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/julian7/utta/internal/remote"
	"github.com/julian7/utta/tunnel"
	"github.com/urfave/cli/v2"
)

func (a *App) remoteCommand() *cli.Command {
	flags := a.CommonFlags()
	flags = append(
		flags,
		&cli.StringFlag{
			Name:     "sshlisten",
			Value:    "",
			Usage:    "SSH remote listening port",
			Required: true,
			EnvVars:  []string{"UTTA_SSH_LISTEN"},
		},
		&cli.StringFlag{
			Name:     "sshconnect",
			Value:    "",
			Usage:    "SSH local target port",
			Required: true,
			EnvVars:  []string{"UTTA_SSH_CONNECT"},
		},
		&cli.UintFlag{
			Name:    "breaker",
			Value:   3,
			Usage:   "Circuit breaker: taking a break after # attempts",
			EnvVars: []string{"UTTA_BREAKER"},
		},
		&cli.DurationFlag{
			Name:    "sleep",
			Value:   30 * time.Minute,
			Usage:   "Sleep between circuit breaks",
			EnvVars: []string{"UTTA_SLEEP"},
		},
	)

	return &cli.Command{
		Name:   "remote",
		Usage:  "create remotely listening tunnel",
		Action: a.remoteAction,
		Flags:  flags,
	}
}

func (a *App) remoteAction(c *cli.Context) error {
	conn, err := a.GenericConnection(c)
	if err != nil {
		return err
	}

	if !conn.HasSSH() {
		return errors.New("no ssh connection configured")
	}

	tun := tunnel.NewConnection(
		a.logger,
		c.String("sshconnect"),
		"",
	)
	if err := tun.SetupTCPConnector(); err != nil {
		return fmt.Errorf("setting up remote connector: %w", err)
	}

	remote.New(conn, a.logger, c.String("sshlisten"), c.Uint("breaker"), c.Duration("sleep")).Loop(tun)

	return nil
}

func (a *App) remoteBefore(c *cli.Context) error {
	if c.String("connect") == "" {
		return errors.New("--connect is required")
	}
	if c.String("sshuser") == "" {
		return errors.New("--sshuser not provided")
	}
	if c.String("sshkey") == "" {
		return errors.New("--sshkey not provided")
	}
	return nil
}
