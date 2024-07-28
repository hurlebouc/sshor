package ssh

/*
references:

	- https://stackoverflow.com/q/58237670/4288267
	- https://github.com/containerd/console/blob/65eb8c0396d0cac15c888bcf4d47049c21317b18/console_windows.go#L191
	- https://learn.microsoft.com/en-us/windows/console/high-level-console-modes
	- https://pkg.go.dev/go/build#hdr-Build_Constraints
*/

import (
	"golang.org/x/crypto/ssh"
	"golang.org/x/sys/windows"
)

type state struct {
	mode uint32
}

func adaptConsole(fd int) error {
	var st uint32
	if err := windows.GetConsoleMode(windows.Handle(fd), &st); err != nil {
		return err
	}
	raw := st | windows.ENABLE_VIRTUAL_TERMINAL_INPUT
	if err := windows.SetConsoleMode(windows.Handle(fd), raw); err != nil {
		return err
	}
	return nil
}

func startWindowChangeListerner(session *ssh.Session, height, width int) {

	// This code is disabled because libraty tsize does not seems to works very well with Windows

	// sizeListener, err := tsize.NewSizeListener()
	// if err != nil {
	// 	panic(err)
	// }

	// go func() {
	// 	currentHeight := height
	// 	currentWidth := width
	// 	for change := range sizeListener.Change {
	// 		if currentHeight != change.Height || currentWidth != change.Width {
	// 			session.WindowChange(change.Height, change.Width)
	// 			currentHeight = change.Height
	// 			currentWidth = change.Width
	// 		}
	// 	}
	// }()
}
