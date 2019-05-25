package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var Server = cli.Command{
	Name:   "Server",
	Usage:  "Start server",
	Action: runServer,
	Flags:  []cli.Flag{},
}

func runServer(ctx *cli.Context) error {
	logrus.Info("RunServer")
	return nil
}
