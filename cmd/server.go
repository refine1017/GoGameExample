package cmd

import (
	"github.com/refine1017/GoGameExample/modules/net"
	"github.com/refine1017/GoGameExample/modules/setting"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"sync"
)

var Server = cli.Command{
	Name:        "Server",
	Description: "Golang Game Server",
	Usage:       "Start server",
	Action:      runServer,
	Flags:       []cli.Flag{},
}

func runServer(ctx *cli.Context) error {
	if err := setting.Load(); err != nil {
		return err
	}

	addr := setting.Server.Host + ":" + setting.Server.Port

	server := net.NewServer(addr)

	waiter := &sync.WaitGroup{}

	if err := server.Run(waiter); err != nil {
		return err
	}

	logrus.Info("Server running...")

	waiter.Wait()

	return nil
}
