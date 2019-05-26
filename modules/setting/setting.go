package setting

import (
	"github.com/Unknwon/com"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

var (
	Cfg        *ini.File
	ConfigFile = "config/app.ini"
)

func Load() error {
	Cfg = ini.Empty()

	if com.IsFile(ConfigFile) {
		if err := Cfg.Append(ConfigFile); err != nil {
			logrus.Fatalf("Failed to load config file '%s': %v", ConfigFile, err)
		}
	} else {
		logrus.Warnf("Config '%s' not found", ConfigFile)
	}

	Cfg.NameMapper = ini.AllCapsUnderscore

	// server
	Server.Host = Cfg.Section("server").Key("Host").MustString("127.0.0.1")
	Server.Port = Cfg.Section("server").Key("Port").MustString("8000")

	return nil
}
