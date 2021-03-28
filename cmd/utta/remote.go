package main

import (
	"errors"
	"fmt"
	"time"

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
	)

	return &cli.Command{
		Name:   "remote",
		Usage:  "create remotely listening tunnel",
		Action: a.remoteAction,
		Flags:  flags,
	}
}

func wait() {
	time.Sleep(5 * time.Second)
}

func (a *App) remoteAction(c *cli.Context) error {
	connection, err := a.GenericConnection(c)
	if err != nil {
		return err
	}

	if !connection.HasSSH() {
		return errors.New("no ssh connection configured")
	}

	intConnection := tunnel.NewConnection(
		a.logger,
		c.String("sshconnect"),
		"",
	)
	if err := intConnection.SetupTCPConnector(); err != nil {
		return fmt.Errorf("setting up remote connector: %w", err)
	}

	for {
		conn, err := connection.Dial()
		if err != nil {
			_ = a.logger.Log(
				"level", "error",
				"msg", "connection error",
				"err", err,
			)
			wait()
			continue
		}

		listener := tunnel.NewRemoteListener(
			c.String("sshlisten"),
			connection,
		)

		err = tunnel.NewTunnel(a.logger, listener).Run(intConnection)
		if err != nil {
			_ = a.logger.Log(
				"level", "error",
				"msg", "SSH tunnel error",
				"err", err,
			)
			_ = conn.Close()
			wait()
			continue
		}
	}
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
