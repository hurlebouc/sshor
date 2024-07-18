package cmd_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hurlebouc/sshor/cmd"
)

func TestCopyDirLocalToRemote(t *testing.T) {
	destDir := initTempDir(directoryLayout{})
	defer func() {
		err := os.RemoveAll(destDir)
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
		dirs: map[string]directoryLayout{
			"emptydir": {},
			"subdir": {
				files: map[string]file{
					"subtest": {
						content: []byte("sub"),
					},
				},
			},
		},
	}
	srcDir := initTempDir(srcLayout)
	defer func() {
		err := os.RemoveAll(srcDir)
		if err != nil {
			panic(err)
		}
	}()

	c := make(chan struct{})
	go startSftpServer(c, "toto", "totopwd", 2344, destDir)
	<-c

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{
		"sshor",
		"copy",
		srcDir,
		"toto@127.0.0.1:",
		"-w",
		"totopwd",
		"-p",
		"2344",
	}
	cmd.Execute()
	copiedLayout := readDirectory(filepath.Join(destDir, filepath.Base(srcDir)))
	if !equalDirs(copiedLayout, srcLayout) {
		t.Fatalf("final directory\n--> %+v\nis distinct from source directory\n--> %+v", copiedLayout, srcLayout)
	}
}

func TestCopyDirDotRemoteToLocal(t *testing.T) {
	destDir := initTempDir(directoryLayout{})
	defer func() {
		err := os.RemoveAll(destDir)
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
		dirs: map[string]directoryLayout{
			"emptydir": {},
			"subdir": {
				files: map[string]file{
					"subtest": {
						content: []byte("sub"),
					},
				},
			},
		},
	}
	srcDir := initTempDir(srcLayout)
	defer func() {
		err := os.RemoveAll(srcDir)
		if err != nil {
			panic(err)
		}
	}()

	c := make(chan struct{})
	go startSftpServer(c, "toto", "totopwd", 2344, srcDir)
	<-c

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{
		"sshor",
		"copy",
		"toto@127.0.0.1:",
		destDir,
		"-w",
		"totopwd",
		"-p",
		"2344",
	}
	cmd.Execute()
	copiedLayout := readDirectory(destDir)
	if !equalDirs(copiedLayout, srcLayout) {
		t.Fatalf("final directory\n--> %+v\nis distinct from source directory\n--> %+v", copiedLayout, srcLayout)
	}
}

func TestCopyDirRemoteToLocal(t *testing.T) {
	destDir := initTempDir(directoryLayout{})
	defer func() {
		err := os.RemoveAll(destDir)
		if err != nil {
			panic(err)
		}
	}()
	srcLayout := directoryLayout{
		dirs: map[string]directoryLayout{
			"sub": {
				files: map[string]file{
					"test1": {
						content: []byte("coucou"),
					},
					"test2": {
						content: []byte("plop"),
					},
					"test3": {},
				},
				dirs: map[string]directoryLayout{
					"emptydir": {},
					"subdir": {
						files: map[string]file{
							"subtest": {
								content: []byte("sub"),
							},
						},
					},
				},
			},
		},
	}
	srcDir := initTempDir(srcLayout)
	defer func() {
		err := os.RemoveAll(srcDir)
		if err != nil {
			panic(err)
		}
	}()

	c := make(chan struct{})
	go startSftpServer(c, "toto", "totopwd", 2344, srcDir)
	<-c

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{
		"sshor",
		"copy",
		"toto@127.0.0.1:sub",
		destDir,
		"-w",
		"totopwd",
		"-p",
		"2344",
	}
	cmd.Execute()
	copiedLayout := readDirectory(destDir)
	if !equalDirs(copiedLayout, srcLayout) {
		t.Fatalf("final directory\n--> %+v\nis distinct from source directory\n--> %+v", copiedLayout, srcLayout)
	}
}
