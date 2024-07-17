package cmd_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hurlebouc/sshor/cmd"
)

func TestCopy(t *testing.T) {
	remoteDir := initTempDir(directoryLayout{})
	defer func() {
		err := os.RemoveAll(remoteDir)
		if err != nil {
			panic(err)
		}
	}()
	srcLayout := directoryLayout{
		files: map[string]file{
			"test1": {
				content: []byte("coucou"),
			},
			"test2": {
				content: []byte("plop"),
			},
			"test3": {},
		},
		dirs: map[string]directoryLayout{},
	}
	localDir := initTempDir(srcLayout)
	defer func() {
		err := os.RemoveAll(localDir)
		if err != nil {
			panic(err)
		}
	}()

	c := make(chan struct{})
	go startSftpServer(c, "toto", "totopwd", 2344, remoteDir)
	<-c

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{
		"sshor",
		"copy",
		localDir,
		"toto@127.0.0.1:",
		"-w",
		"totopwd",
		"-p",
		"2344",
	}
	cmd.Execute()
	newRemoteLayout := readDirectory(filepath.Join(remoteDir, filepath.Base(localDir)))
	if !equalDirs(newRemoteLayout, srcLayout) {
		t.Fatalf("remote directory\n--> %+v\nis distinct from local directory\n--> %+v", newRemoteLayout, srcLayout)
	}
}
