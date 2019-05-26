package main

import (
	"github.com/refine1017/GoGameExample/cmd"
	"github.com/refine1017/GoGameExample/modules/setting"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "GoGameExample"
	app.Description = "Golang Game Example"
	app.Version = "1.0.0"
	app.Commands = []cli.Command{
		cmd.Server,
	}

	defaultFlags := []cli.Flag{
		cli.StringFlag{
			Name:        "config, c",
			Value:       setting.ConfigFile,
			Usage:       "config file path",
			Destination: &setting.ConfigFile,
		},
		cli.VersionFlag,
	}

	app.Flags = append(app.Flags, cmd.Server.Flags...)
	app.Flags = append(app.Flags, defaultFlags...)
	app.Action = cmd.Server.Action

	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatalf("Failed to run app with %s: %v", os.Args, err)
	}
}
