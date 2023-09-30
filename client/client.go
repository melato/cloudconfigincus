package client

import (
	"fmt"

	incus "github.com/lxc/incus/client"
	config "github.com/lxc/incus/shared/cliconfig"
)

var Trace bool

/*
InstanceClient provides an InstanceServer
Use it by calling these methods, in order:
Configured(), CurrentServer().
*/
type InstanceClient struct {
	ForceLocal bool   `name:"force-local" usage:"Force using the local unix socket"`
	Project    string `name:"project" usage:"Override the default project"`

	conf       *config.Config
	rootServer incus.InstanceServer
}

func (t *InstanceClient) Init() error {
	return nil
}

func connectUnix() (incus.InstanceServer, error) {
	unixSocket, err := UnixSocket()
	if err != nil {
		return nil, err
	}
	if unixSocket == "" || !PathExists(unixSocket) {
		return nil, fmt.Errorf("no such unix socket: %s", unixSocket)
	}
	if Trace {
		fmt.Printf("using unix socket: %s\n", unixSocket)
	}
	server, err := incus.ConnectIncusUnix(unixSocket, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", unixSocket, err)
	}
	return server, nil
}

func (t *InstanceClient) Configured() error {
	if t.ForceLocal {
		server, err := connectUnix()
		if err != nil {
			return err
		}
		t.rootServer = server
		t.conf = config.NewConfig("", true)
	} else {
		conf, err := LoadConfig()
		if err != nil {
			return err
		}
		t.conf = conf
	}
	return nil
}

// RootServer - return the unqualified (no project) instance server
func (t *InstanceClient) RootServer() (incus.InstanceServer, error) {
	if t.rootServer == nil {
		remote, ok := t.conf.Remotes[t.conf.DefaultRemote]
		if ok && remote.Addr == "unix://" {
			server, err := connectUnix()
			if err != nil {
				return nil, err
			}
			t.rootServer = server
		} else {
			d, err := t.conf.GetInstanceServer(t.conf.DefaultRemote)
			if err != nil {
				return nil, err
			}
			t.rootServer = d
		}
	}
	return t.rootServer, nil
}

// RootServer - return the instance server for the specified project
// If project is empty, use the default project
func (t *InstanceClient) ProjectServer(project string) (incus.InstanceServer, error) {
	if project == "" {
		project = t.CurrentProject()
	}

	server, err := t.RootServer()
	if err != nil {
		return nil, err
	}
	return server.UseProject(project), nil
}

// RootServer - return the instance server for the current project
func (t *InstanceClient) CurrentServer() (incus.InstanceServer, error) {
	return t.ProjectServer("")
}

func (t *InstanceClient) CurrentProject() string {
	if t.Project != "" {
		return t.Project
	}
	if t.conf != nil {
		remote, exists := t.conf.Remotes[t.conf.DefaultRemote]
		if exists {
			return remote.Project
		}
	}
	return ""
}

func (t *InstanceClient) Config() *config.Config {
	return t.conf
}
