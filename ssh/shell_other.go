//go:build !windows

package ssh

import (
	tsize "github.com/kopoli/go-terminal-size"
	"golang.org/x/crypto/ssh"
)

func adaptConsole(_ int) error {
	return nil
}

func startWindowChangeListerner(session *ssh.Session, height, width int) {
	sizeListener, err := tsize.NewSizeListener()
	if err != nil {
		panic(err)
	}

	go func() {
		currentHeight := height
		currentWidth := width
		for change := range sizeListener.Change {
			if currentHeight != change.Height || currentWidth != change.Width {
				session.WindowChange(change.Height, change.Width)
				currentHeight = change.Height
				currentWidth = change.Width
			}
		}
	}()
}
