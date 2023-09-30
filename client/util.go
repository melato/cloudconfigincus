package client

import (
	"fmt"
	"os"
	"path/filepath"

	config "github.com/lxc/incus/shared/cliconfig"
)

func PathExists(name string) bool {
	_, err := os.Stat(name)
	if err != nil && os.IsNotExist(err) {
		return false
	}

	return true
}

/*
UnixSocket returns the default Unix Socket
See https://github.com/lxc/incus/blob/main/client/connection.go
We use a similar strategy as connection.go ConnectIncusUnixWithContext(),
*/
func UnixSocket() (string, error) {
	path := os.Getenv("INCUS_SOCKET")
	if path == "" {
		incusDir := os.Getenv("INCUS_DIR")
		if incusDir == "" {
			incusDir = "/var/lib/incus"
		}

		path = filepath.Join(incusDir, "unix.socket")
	}
	return path, nil

	//return "", fmt.Errorf("no incus socket found")
}

/*
ConfigDir returns the default incus client configuration directory.
See https://github.com/lxc/incus/blob/master/cmd/incus/main.go
*/
func ConfigDir() (string, error) {
	configDir := os.Getenv("INCUS_CONF")
	if configDir != "" {
		return configDir, nil
	}

	userConfigDir, err := os.UserConfigDir()
	if err == nil && userConfigDir != "" {
		configDir = filepath.Join(userConfigDir, "incus")
		if PathExists(configDir) {
			return configDir, nil
		}
	}
	var c config.Config
	configDir = c.GlobalConfigPath()
	if PathExists(configDir) {
		return configDir, nil
	}
	return "", nil
}

func LoadConfig() (*config.Config, error) {
	configDir, err := ConfigDir()
	if err != nil {
		return nil, err
	}
	if Trace {
		fmt.Printf("using %s\n", configDir)
	}
	confPath := os.ExpandEnv(filepath.Join(configDir, "config.yml"))
	return config.LoadConfig(confPath)
}
