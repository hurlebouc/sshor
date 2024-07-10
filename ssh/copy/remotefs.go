package copy

import (
	"io/fs"

	"github.com/hurlebouc/sshor/config"
	"github.com/hurlebouc/sshor/ssh"
	"github.com/pkg/sftp"
)

type remoteFS struct {
	client    *sftp.Client
	sshClient ssh.SshClient
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

func (remote remoteFS) close() {
	remote.client.Close()
	remote.sshClient.Close()
}

func NewRemote(hostConfig config.Host, passwordFlag, keepassPwdFlag string, path string) Endpoint {
	keepassPwdMap := ssh.InitKeepassPwdMap(hostConfig, keepassPwdFlag)

	ctx := ssh.InitContext()
	sshClient, _ := ssh.NewSshClient(ctx, hostConfig, passwordFlag, keepassPwdMap)
	if sshClient.Client == nil {
		panic("Cannot construct ssh client. This is probably caused by not specifying host of the target.")
	}
	return Endpoint{
		path: path,
		fileSystem: remoteFS{
			client:    newSftp(sshClient.Client),
			sshClient: sshClient,
		},
	}
}

func NewRemoteFromFS(remoteFS remoteFS, path string) Endpoint {
	return Endpoint{
		path:       path,
		fileSystem: remoteFS,
	}
}
