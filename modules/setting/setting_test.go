package setting

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadNoFile(t *testing.T) {
	if err := Load(); err != nil {
		t.Errorf("Load with err: %v", err)
	}

	assert.Equal(t, "127.0.0.1", Server.Host, "Server Host Error")
	assert.Equal(t, "8000", Server.Port, "Server Port Error")
}

func TestLoadFile(t *testing.T) {
	ConfigFile = "../../config/app.ini"

	if err := Load(); err != nil {
		t.Errorf("Load %v with err: %v", ConfigFile, err)
	}

	assert.Equal(t, "localhost", Server.Host, "Server Host Error")
	assert.Equal(t, "8800", Server.Port, "Server Port Error")
}

func TestLoadFileNoExist(t *testing.T) {
	ConfigFile = "../config/no.ini"

	if err := Load(); err != nil {
		t.Errorf("Load %v with err: %v", ConfigFile, err)
	}
}
