package cmd

import (
	"os"
	"testing"
)

func TestCopy(t *testing.T) {
	dirName := initTempDir(directoryLayout{})
	defer func() {
		err := os.RemoveAll(dirName)
		if err != nil {
			panic(err)
		}
	}()
	c := make(chan struct{})
	go startSftpServer(c, "toto", "totopwd", 2344, dirName)
	<-c
}
