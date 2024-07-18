package cmd_test

import (
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
	go startSftpServer(c, "toto", "totopwd", 2344, destDir)
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
		"2344",
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
	go startSftpServer(c, "toto", "totopwd", 2344, destDir)
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
		"2344",
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
	go startSftpServer(c, "toto", "totopwd", 2344, destDir)
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
		"2344",
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
	if !equalDirs(copiedLayout, expectedLayout) {
		t.Fatalf("final directory\n--> %+v\nis distinct from source directory\n--> %+v", copiedLayout, srcLayout)
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
	go startSftpServer(c, "toto", "totopwd", 2344, srcDir)
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
		"2344",
	}
	cmd.Execute()
	copiedLayout := readDirectory(destDir)
	if !equalDirs(copiedLayout, expectedLayout) {
		t.Fatalf("final directory\n--> %+v\nis distinct from source directory\n--> %+v", copiedLayout, srcLayout)
	}
}
