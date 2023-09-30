package main

import (
	_ "embed"
	"fmt"
	"os"

	incus "github.com/lxc/incus/client"
	"melato.org/cloudconfig"
	"melato.org/cloudconfig/ostype"
	"melato.org/cloudconfigincus"
	"melato.org/cloudconfigincus/client"
	"melato.org/command"
	"melato.org/command/usage"
)

//go:embed version
var version string

//go:embed usage.yaml
var usageData []byte

type App struct {
	Client   client.InstanceClient
	Instance string `name:"i" usage:"instance to configure"`
	OS       string `name:"ostype" usage:"OS type"`
	os       cloudconfig.OSType
	server   incus.InstanceServer
}

func (t *App) Init() error {
	return t.Client.Init()
}

func (t *App) Configured() error {
	if t.Instance == "" {
		return fmt.Errorf("missing instance")
	}
	err := t.Client.Configured()
	if err != nil {
		return err
	}
	server, err := t.Client.CurrentServer()
	if err != nil {
		return err
	}
	switch t.OS {
	case "":
	case "alpine":
		t.os = &ostype.Alpine{}
	case "debian":
		t.os = &ostype.Debian{}
	default:
		return fmt.Errorf("unknown OS type: %s", t.OS)
	}
	t.server = server
	return nil
}

func (t *App) Apply(configFiles ...string) error {
	base := cloudconfigincus.NewInstanceConfigurer(t.server, t.Instance)
	base.Log = os.Stdout
	configurer := cloudconfig.NewConfigurer(base)
	configurer.OS = t.os
	configurer.Log = os.Stdout
	if len(configFiles) == 1 && configFiles[0] == "-" {
		return configurer.ApplyStdin()
	} else {
		return configurer.ApplyConfigFiles(configFiles...)
	}
}

func (t *App) FileExists(path string) error {
	base := cloudconfigincus.NewInstanceConfigurer(t.server, t.Instance)
	base.Log = os.Stdout
	exists, err := base.FileExists(path)
	if err != nil {
		return err
	}
	fmt.Printf("%s: %v\n", path, exists)
	return nil
}

func main() {
	cmd := &command.SimpleCommand{}
	var app App
	cmd.Command("apply").Flags(&app).RunFunc(app.Apply)
	cmd.Command("file-exists").Flags(&app).RunFunc(app.FileExists)
	cmd.Command("version").RunFunc(func() { fmt.Println(version) })

	usage.Apply(cmd, usageData)
	command.Main(cmd)
}
