package cmd_test

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/hurlebouc/sshor/cmd"
)

func TestCopyDirDotLocalToRemote(t *testing.T) {
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
	port := uint16(rand.Uint32())
	if port <= 1024 {
		port = port + 1024
	}
	go startSftpServer(c, "toto", "totopwd", port, destDir)
	<-c

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{
		"sshor",
		"copy",
		".",
		"toto@127.0.0.1:",
		"-w",
		"totopwd",
		"-p",
		fmt.Sprintf("%d", port),
	}
	oldwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	os.Chdir(srcDir)
	defer os.Chdir(oldwd)

	cmd.Execute()
	copiedLayout := readDirectory(destDir)
	if !equalDirs(copiedLayout, srcLayout) {
		t.Fatalf("final directory\n--> %+v\nis distinct from source directory\n--> %+v", copiedLayout, srcLayout)
	}
}

func TestCopyDirLocalToDotRemote(t *testing.T) {
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
	port := uint16(rand.Uint32())
	if port <= 1024 {
		port = port + 1024
	}
	go startSftpServer(c, "toto", "totopwd", port, destDir)
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
		fmt.Sprintf("%d", port),
	}
	cmd.Execute()
	copiedLayout := readDirectory(filepath.Join(destDir, filepath.Base(srcDir)))
	if !equalDirs(copiedLayout, srcLayout) {
		t.Fatalf("final directory\n--> %+v\nis distinct from source directory\n--> %+v", copiedLayout, srcLayout)
	}
}

func TestCopyDirMissingLocalToRemote(t *testing.T) {
	destDir := initTempDir(directoryLayout{
		dirs: map[string]directoryLayout{
			"sub": {},
		},
	})
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
	port := uint16(rand.Uint32())
	if port <= 1024 {
		port = port + 1024
	}
	go startSftpServer(c, "toto", "totopwd", port, destDir)
	<-c

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{
		"sshor",
		"copy",
		filepath.Join(srcDir, "plop"),
		"toto@127.0.0.1:sub",
		"-w",
		"totopwd",
		"-p",
		fmt.Sprintf("%d", port),
	}
	expectPanic(t, cmd.Execute)
}
func TestCopyDirLocalToRemote(t *testing.T) {
	destDir := initTempDir(directoryLayout{
		dirs: map[string]directoryLayout{
			"sub": {},
		},
	})
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
	port := uint16(rand.Uint32())
	if port <= 1024 {
		port = port + 1024
	}
	go startSftpServer(c, "toto", "totopwd", port, destDir)
	<-c

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{
		"sshor",
		"copy",
		srcDir,
		"toto@127.0.0.1:sub",
		"-w",
		"totopwd",
		"-p",
		fmt.Sprintf("%d", port),
	}
	cmd.Execute()
	copiedLayout := readDirectory(filepath.Join(destDir, "sub", filepath.Base(srcDir)))
	if !equalDirs(copiedLayout, srcLayout) {
		t.Fatalf("final directory\n--> %+v\nis distinct from source directory\n--> %+v", copiedLayout, srcLayout)
	}
}

func TestCopyDirLocalToMissingRemote(t *testing.T) {
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
	port := uint16(rand.Uint32())
	if port <= 1024 {
		port = port + 1024
	}
	go startSftpServer(c, "toto", "totopwd", port, destDir)
	<-c

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{
		"sshor",
		"copy",
		srcDir,
		"toto@127.0.0.1:sub",
		"-w",
		"totopwd",
		"-p",
		fmt.Sprintf("%d", port),
	}
	cmd.Execute()
	copiedLayout := readDirectory(filepath.Join(destDir, "sub"))
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
	port := uint16(rand.Uint32())
	if port <= 1024 {
		port = port + 1024
	}
	go startSftpServer(c, "toto", "totopwd", port, srcDir)
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
		fmt.Sprintf("%d", port),
	}
	cmd.Execute()
	copiedLayout := readDirectory(destDir)
	if !equalDirs(copiedLayout, srcLayout) {
		t.Fatalf("final directory\n--> %+v\nis distinct from source directory\n--> %+v", copiedLayout, srcLayout)
	}
}

func TestCopyDirMissingRemoteToLocal(t *testing.T) {
	destDir := initTempDir(directoryLayout{})
	defer func() {
		err := os.RemoveAll(destDir)
		if err != nil {
			panic(err)
		}
	}()
	subdir := directoryLayout{
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
	srcLayout := directoryLayout{
		files: map[string]file{
			"foo": {
				content: []byte("it's a trap!"),
			},
		},
		dirs: map[string]directoryLayout{
			"sub": subdir,
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
	port := uint16(rand.Uint32())
	if port <= 1024 {
		port = port + 1024
	}
	go startSftpServer(c, "toto", "totopwd", port, srcDir)
	<-c

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{
		"sshor",
		"copy",
		"toto@127.0.0.1:plop",
		destDir,
		"-w",
		"totopwd",
		"-p",
		fmt.Sprintf("%d", port),
	}
	expectPanic(t, cmd.Execute)
}
func TestCopyDirRemoteToLocal(t *testing.T) {
	destDir := initTempDir(directoryLayout{})
	defer func() {
		err := os.RemoveAll(destDir)
		if err != nil {
			panic(err)
		}
	}()
	subdir := directoryLayout{
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
	srcLayout := directoryLayout{
		files: map[string]file{
			"foo": {
				content: []byte("it's a trap!"),
			},
		},
		dirs: map[string]directoryLayout{
			"sub": subdir,
		},
	}
	srcDir := initTempDir(srcLayout)
	defer func() {
		err := os.RemoveAll(srcDir)
		if err != nil {
			panic(err)
		}
	}()
	expectedLayout := directoryLayout{
		dirs: map[string]directoryLayout{
			"sub": subdir,
		},
	}

	c := make(chan struct{})
	port := uint16(rand.Uint32())
	if port <= 1024 {
		port = port + 1024
	}
	go startSftpServer(c, "toto", "totopwd", port, srcDir)
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
		fmt.Sprintf("%d", port),
	}
	cmd.Execute()
	copiedLayout := readDirectory(destDir)
	if !equalDirs(copiedLayout, expectedLayout) {
		t.Fatalf("final directory\n--> %+v\nis distinct from source directory\n--> %+v", copiedLayout, expectedLayout)
	}
}

func TestCopyDirRemoteToMissingLocal(t *testing.T) {
	destDir := initTempDir(directoryLayout{})
	defer func() {
		err := os.RemoveAll(destDir)
		if err != nil {
			panic(err)
		}
	}()
	subdir := directoryLayout{
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
	srcLayout := directoryLayout{
		files: map[string]file{
			"foo": {
				content: []byte("it's a trap!"),
			},
		},
		dirs: map[string]directoryLayout{
			"sub": subdir,
		},
	}
	srcDir := initTempDir(srcLayout)
	defer func() {
		err := os.RemoveAll(srcDir)
		if err != nil {
			panic(err)
		}
	}()
	expectedLayout := directoryLayout{
		dirs: map[string]directoryLayout{
			"plop": subdir,
		},
	}

	c := make(chan struct{})
	port := uint16(rand.Uint32())
	if port <= 1024 {
		port = port + 1024
	}
	go startSftpServer(c, "toto", "totopwd", port, srcDir)
	<-c

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{
		"sshor",
		"copy",
		"toto@127.0.0.1:sub",
		filepath.Join(destDir, "plop"),
		"-w",
		"totopwd",
		"-p",
		fmt.Sprintf("%d", port),
	}
	cmd.Execute()
	copiedLayout := readDirectory(destDir)
	if !equalDirs(copiedLayout, expectedLayout) {
		t.Fatalf("final directory\n--> %+v\nis distinct from source directory\n--> %+v", copiedLayout, expectedLayout)
	}
}

func TestCopyFileRemoteToLocal(t *testing.T) {
	destDir := initTempDir(directoryLayout{})
	defer func() {
		err := os.RemoveAll(destDir)
		if err != nil {
			panic(err)
		}
	}()
	srcLayout := directoryLayout{
		files: map[string]file{
			"plop": {
				content: []byte("coucou"),
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
	port := uint16(rand.Uint32())
	if port <= 1024 {
		port = port + 1024
	}
	go startSftpServer(c, "toto", "totopwd", port, srcDir)
	<-c

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{
		"sshor",
		"copy",
		"toto@127.0.0.1:plop",
		destDir,
		"-w",
		"totopwd",
		"-p",
		fmt.Sprintf("%d", port),
	}
	cmd.Execute()
	copiedLayout := readDirectory(destDir)
	if !equalDirs(copiedLayout, srcLayout) {
		t.Fatalf("final directory\n--> %+v\nis distinct from source directory\n--> %+v", copiedLayout, srcLayout)
	}
}

func TestCopyFileRemoteToFileLocal(t *testing.T) {
	destDir := initTempDir(directoryLayout{})
	defer func() {
		err := os.RemoveAll(destDir)
		if err != nil {
			panic(err)
		}
	}()
	srcLayout := directoryLayout{
		files: map[string]file{
			"plop": {
				content: []byte("coucou"),
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
	expectedLayout := directoryLayout{
		files: map[string]file{
			"plip": {
				content: []byte("coucou"),
			},
		},
	}

	c := make(chan struct{})
	port := uint16(rand.Uint32())
	if port <= 1024 {
		port = port + 1024
	}
	go startSftpServer(c, "toto", "totopwd", port, srcDir)
	<-c

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{
		"sshor",
		"copy",
		"toto@127.0.0.1:plop",
		filepath.Join(destDir, "plip"),
		"-w",
		"totopwd",
		"-p",
		fmt.Sprintf("%d", port),
	}
	cmd.Execute()
	copiedLayout := readDirectory(destDir)
	if !equalDirs(copiedLayout, expectedLayout) {
		t.Fatalf("final directory\n--> %+v\nis distinct from source directory\n--> %+v", copiedLayout, expectedLayout)
	}
}

func TestCopyFileLocalToRemote(t *testing.T) {
	destDir := initTempDir(directoryLayout{})
	defer func() {
		err := os.RemoveAll(destDir)
		if err != nil {
			panic(err)
		}
	}()
	srcLayout := directoryLayout{
		files: map[string]file{
			"plop": {
				content: []byte("coucou"),
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
	port := uint16(rand.Uint32())
	if port <= 1024 {
		port = port + 1024
	}
	go startSftpServer(c, "toto", "totopwd", port, destDir)
	<-c

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{
		"sshor",
		"copy",
		filepath.Join(srcDir, "plop"),
		"toto@127.0.0.1:",
		"-w",
		"totopwd",
		"-p",
		fmt.Sprintf("%d", port),
	}
	cmd.Execute()
	copiedLayout := readDirectory(destDir)
	if !equalDirs(copiedLayout, srcLayout) {
		t.Fatalf("final directory\n--> %+v\nis distinct from source directory\n--> %+v", copiedLayout, srcLayout)
	}
}

func TestCopyFileLocalToFileRemote(t *testing.T) {
	destDir := initTempDir(directoryLayout{})
	defer func() {
		err := os.RemoveAll(destDir)
		if err != nil {
			panic(err)
		}
	}()
	srcLayout := directoryLayout{
		files: map[string]file{
			"plop": {
				content: []byte("coucou"),
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
	expectedLayout := directoryLayout{
		files: map[string]file{
			"plip": {
				content: []byte("coucou"),
			},
		},
	}

	c := make(chan struct{})
	port := uint16(rand.Uint32())
	if port <= 1024 {
		port = port + 1024
	}
	go startSftpServer(c, "toto", "totopwd", port, destDir)
	<-c

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{
		"sshor",
		"copy",
		filepath.Join(srcDir, "plop"),
		"toto@127.0.0.1:plip",
		"-w",
		"totopwd",
		"-p",
		fmt.Sprintf("%d", port),
	}
	cmd.Execute()
	copiedLayout := readDirectory(destDir)
	if !equalDirs(copiedLayout, expectedLayout) {
		t.Fatalf("final directory\n--> %+v\nis distinct from source directory\n--> %+v", copiedLayout, expectedLayout)
	}
}

func TestCopyFileMissingRemoteToLocal(t *testing.T) {
	destDir := initTempDir(directoryLayout{})
	defer func() {
		err := os.RemoveAll(destDir)
		if err != nil {
			panic(err)
		}
	}()
	srcLayout := directoryLayout{
		files: map[string]file{
			"plop": {
				content: []byte("coucou"),
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
	port := uint16(rand.Uint32())
	if port <= 1024 {
		port = port + 1024
	}
	go startSftpServer(c, "toto", "totopwd", port, srcDir)
	<-c

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{
		"sshor",
		"copy",
		"toto@127.0.0.1:plip",
		destDir,
		"-w",
		"totopwd",
		"-p",
		fmt.Sprintf("%d", port),
	}
	expectPanic(t, cmd.Execute)
}

func TestCopyFileMissingLocalToRemote(t *testing.T) {
	destDir := initTempDir(directoryLayout{})
	defer func() {
		err := os.RemoveAll(destDir)
		if err != nil {
			panic(err)
		}
	}()
	srcLayout := directoryLayout{
		files: map[string]file{
			"plop": {
				content: []byte("coucou"),
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
	port := uint16(rand.Uint32())
	if port <= 1024 {
		port = port + 1024
	}
	go startSftpServer(c, "toto", "totopwd", port, destDir)
	<-c

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{
		"sshor",
		"copy",
		filepath.Join(srcDir, "plip"),
		"toto@127.0.0.1:",
		"-w",
		"totopwd",
		"-p",
		fmt.Sprintf("%d", port),
	}
	expectPanic(t, cmd.Execute)
}

func TestCopyFileLocalToSubRemote(t *testing.T) {
	destDir := initTempDir(directoryLayout{})
	defer func() {
		err := os.RemoveAll(destDir)
		if err != nil {
			panic(err)
		}
	}()
	srcLayout := directoryLayout{
		files: map[string]file{
			"plop": {
				content: []byte("coucou"),
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
	expectedLayout := directoryLayout{
		dirs: map[string]directoryLayout{
			"plif": {
				files: map[string]file{
					"plaf": {
						content: []byte("coucou"),
					},
				},
			},
		},
	}

	c := make(chan struct{})
	port := uint16(rand.Uint32())
	if port <= 1024 {
		port = port + 1024
	}
	go startSftpServer(c, "toto", "totopwd", port, destDir)
	<-c

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{
		"sshor",
		"copy",
		filepath.Join(srcDir, "plop"),
		"toto@127.0.0.1:plif/plaf",
		"-w",
		"totopwd",
		"-p",
		fmt.Sprintf("%d", port),
	}
	cmd.Execute()
	copiedLayout := readDirectory(destDir)
	if !equalDirs(copiedLayout, expectedLayout) {
		t.Fatalf("final directory\n--> %+v\nis distinct from source directory\n--> %+v", copiedLayout, expectedLayout)
	}
}

func TestCopyFileRemoteToSubLocal(t *testing.T) {
	destDir := initTempDir(directoryLayout{})
	defer func() {
		err := os.RemoveAll(destDir)
		if err != nil {
			panic(err)
		}
	}()
	srcLayout := directoryLayout{
		files: map[string]file{
			"plop": {
				content: []byte("coucou"),
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
	expectedLayout := directoryLayout{
		dirs: map[string]directoryLayout{
			"plif": {
				files: map[string]file{
					"plaf": {
						content: []byte("coucou"),
					},
				},
			},
		},
	}

	c := make(chan struct{})
	port := uint16(rand.Uint32())
	if port <= 1024 {
		port = port + 1024
	}
	go startSftpServer(c, "toto", "totopwd", port, srcDir)
	<-c

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{
		"sshor",
		"copy",
		"toto@127.0.0.1:plop",
		filepath.Join(destDir, "plif", "plaf"),
		"-w",
		"totopwd",
		"-p",
		fmt.Sprintf("%d", port),
	}
	cmd.Execute()
	copiedLayout := readDirectory(destDir)
	if !equalDirs(copiedLayout, expectedLayout) {
		t.Fatalf("final directory\n--> %+v\nis distinct from source directory\n--> %+v", copiedLayout, expectedLayout)
	}
}

func TestCopyFileLocalToSubFileRemote(t *testing.T) {
	destDir := initTempDir(directoryLayout{
		files: map[string]file{
			"plif": {},
		},
	})
	defer func() {
		err := os.RemoveAll(destDir)
		if err != nil {
			panic(err)
		}
	}()
	srcLayout := directoryLayout{
		files: map[string]file{
			"plop": {
				content: []byte("coucou"),
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
	port := uint16(rand.Uint32())
	if port <= 1024 {
		port = port + 1024
	}
	go startSftpServer(c, "toto", "totopwd", port, destDir)
	<-c

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{
		"sshor",
		"copy",
		filepath.Join(srcDir, "plop"),
		"toto@127.0.0.1:plif/plaf",
		"-w",
		"totopwd",
		"-p",
		fmt.Sprintf("%d", port),
	}
	expectPanic(t, cmd.Execute)
}

func TestCopyFileRemoteToSubFileLocal(t *testing.T) {
	destDir := initTempDir(directoryLayout{
		files: map[string]file{
			"plif": {},
		},
	})
	defer func() {
		err := os.RemoveAll(destDir)
		if err != nil {
			panic(err)
		}
	}()
	srcLayout := directoryLayout{
		files: map[string]file{
			"plop": {
				content: []byte("coucou"),
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
	port := uint16(rand.Uint32())
	if port <= 1024 {
		port = port + 1024
	}
	go startSftpServer(c, "toto", "totopwd", port, srcDir)
	<-c

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{
		"sshor",
		"copy",
		"toto@127.0.0.1:plop",
		filepath.Join(destDir, "plif", "plaf"),
		"-w",
		"totopwd",
		"-p",
		fmt.Sprintf("%d", port),
	}
	expectPanic(t, cmd.Execute)
}

func TestCopyFileLocalToExistingRemote(t *testing.T) {
	destDir := initTempDir(directoryLayout{
		files: map[string]file{
			"plop": {
				content: []byte("autre"),
			},
		},
	})
	defer func() {
		err := os.RemoveAll(destDir)
		if err != nil {
			panic(err)
		}
	}()
	srcLayout := directoryLayout{
		files: map[string]file{
			"plop": {
				content: []byte("coucou"),
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
	port := uint16(rand.Uint32())
	if port <= 1024 {
		port = port + 1024
	}
	go startSftpServer(c, "toto", "totopwd", port, destDir)
	<-c

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{
		"sshor",
		"copy",
		filepath.Join(srcDir, "plop"),
		"toto@127.0.0.1:",
		"-w",
		"totopwd",
		"-p",
		fmt.Sprintf("%d", port),
	}
	cmd.Execute()
	copiedLayout := readDirectory(destDir)
	if !equalDirs(copiedLayout, srcLayout) {
		t.Fatalf("final directory\n--> %+v\nis distinct from source directory\n--> %+v", copiedLayout, srcLayout)
	}
}

func TestCopyFileRemoteToExistingLocal(t *testing.T) {
	destDir := initTempDir(directoryLayout{
		files: map[string]file{
			"plop": {
				content: []byte("autre"),
			},
		},
	})
	defer func() {
		err := os.RemoveAll(destDir)
		if err != nil {
			panic(err)
		}
	}()
	srcLayout := directoryLayout{
		files: map[string]file{
			"plop": {
				content: []byte("coucou"),
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
	port := uint16(rand.Uint32())
	if port <= 1024 {
		port = port + 1024
	}
	go startSftpServer(c, "toto", "totopwd", port, srcDir)
	<-c

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{
		"sshor",
		"copy",
		"toto@127.0.0.1:plop",
		destDir,
		"-w",
		"totopwd",
		"-p",
		fmt.Sprintf("%d", port),
	}
	cmd.Execute()
	copiedLayout := readDirectory(destDir)
	if !equalDirs(copiedLayout, srcLayout) {
		t.Fatalf("final directory\n--> %+v\nis distinct from source directory\n--> %+v", copiedLayout, srcLayout)
	}
}

func TestCopyDirLocalToFileRemote(t *testing.T) {
	destDir := initTempDir(directoryLayout{
		files: map[string]file{
			"sub": {},
		},
	})
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
	port := uint16(rand.Uint32())
	if port <= 1024 {
		port = port + 1024
	}
	go startSftpServer(c, "toto", "totopwd", port, destDir)
	<-c

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{
		"sshor",
		"copy",
		srcDir,
		"toto@127.0.0.1:sub",
		"-w",
		"totopwd",
		"-p",
		fmt.Sprintf("%d", port),
	}
	expectPanic(t, cmd.Execute)
}

func TestCopyDirRemoteToFileLocal(t *testing.T) {
	destDir := initTempDir(directoryLayout{
		files: map[string]file{
			"sub": {},
		},
	})
	defer func() {
		err := os.RemoveAll(destDir)
		if err != nil {
			panic(err)
		}
	}()
	subdir := directoryLayout{
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
	srcLayout := directoryLayout{
		files: map[string]file{
			"foo": {
				content: []byte("it's a trap!"),
			},
		},
		dirs: map[string]directoryLayout{
			"sub": subdir,
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
	port := uint16(rand.Uint32())
	if port <= 1024 {
		port = port + 1024
	}
	go startSftpServer(c, "toto", "totopwd", port, srcDir)
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
		fmt.Sprintf("%d", port),
	}
	expectPanic(t, cmd.Execute)
}
