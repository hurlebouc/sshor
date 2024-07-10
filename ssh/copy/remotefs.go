package copy

import (
	"context"
	"io/fs"

	"github.com/hurlebouc/sshor/config"
	"github.com/hurlebouc/sshor/ssh"
	"github.com/pkg/sftp"
)

type remoteFS struct {
	client *sftp.Client
}

func (remote remoteFS) join(a, b string) string {
	return a + "/" + b
}

func (remote remoteFS) mkdirAll(path string) error {
	return remote.client.MkdirAll(path)
}

func (remote remoteFS) create(path string) (writer, error) {
	return remote.client.Create(path)
}

func (remote remoteFS) open(path string) (reader, error) {
	return remote.client.Open(path)
}

func (remote remoteFS) readDir(path string) ([]fs.FileInfo, error) {
	return remote.client.ReadDir(path)
}

func NewRemote(ctx context.Context, hostConfig config.Host, passwordFlag string, keepassPwdMap map[string]string, path string) Endpoint {
	sshClient, _ := ssh.NewSshClient(ctx, hostConfig, passwordFlag, keepassPwdMap)
	if sshClient.Client == nil {
		panic("Cannot construct ssh client. This is probably caused by not specifying host of the target.")
	}
	return Endpoint{
		path: path,
		fileSystem: remoteFS{
			client: newSftp(sshClient.Client),
		},
	}
}
