package ssh

import (
	"bytes"
	"fmt"
	"os"

	"github.com/hurlebouc/sshor/config"
	tsize "github.com/kopoli/go-terminal-size"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

type PatternDetector struct {
	buffer   []byte
	position int // next byte to write
}

func (p *PatternDetector) addBytes(b []byte) {
	for i, v := range b {
		p.buffer[(i+p.position)%len(p.buffer)] = v
	}
	p.position = (p.position + len(b)) % len(p.buffer)
}

func (p PatternDetector) exportBytes() []byte {
	res := make([]byte, len(p.buffer))
	for i := range res {
		res[i] = p.buffer[(i+p.position)%len(p.buffer)]
	}
	return res
}

func Shell(hostConf config.Host, passwordFlag, keepassPwdFlag string) {
	keepassPwdMap := InitKeepassPwdMap(hostConf, keepassPwdFlag)

	ctx := InitContext()

	conn, _ := NewSshClient(ctx, hostConf, passwordFlag, keepassPwdMap)
	defer conn.Close()
	// Create a session
	sshClient := GetFirstNonNilSshClient(conn)
	if sshClient == nil {
		panic("todo")
	}
	session, err := sshClient.NewSession()
	if err != nil {
		panic(err)
	}
	defer session.Close()
	// Set up terminal modes
	modes := ssh.TerminalModes{
		//ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		//ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		panic(err)
	}
	defer term.Restore(fd, oldState)

	_ = adaptConsole(fd) // do nothing in case of error

	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		panic(err)
	}
	// Request pseudo terminal
	if err := session.RequestPty("xterm-256color", height, width, modes); err != nil {
		panic(err)
	}

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

	//session.Stdin = os.Stdin
	input, err := session.StdinPipe()
	if err != nil {
		panic(err)
	}
	//session.Stdout = os.Stdout
	output, err := session.StdoutPipe()
	if err != nil {
		panic(err)
	}
	cha := make(chan bool)
	session.Stderr = os.Stderr

	go func() {
		if conn.ChangeUser != nil {
			// se connecte uniquement avec le dernier utilisateur
			_, err := input.Write([]byte(fmt.Sprintf("su - %s\n", conn.ChangeUser.login)))
			if err != nil {
				panic(err)
			}
			<-cha
			input.Write([]byte(fmt.Sprintf("%s\n", conn.ChangeUser.password)))
		}
		buffer := make([]byte, 5)
		for {
			n, err := os.Stdin.Read(buffer)
			if n == 0 {
				break
			}
			if err != nil {
				panic(err)
			}
			_, err = input.Write(buffer[0:n])
			if err != nil {
				panic(err)
			}
		}
	}()

	go func() {
		buffer := make([]byte, 5)
		p := PatternDetector{
			buffer: make([]byte, 100),
		}
		passed := conn.ChangeUser == nil
		for {
			n, err := output.Read(buffer)
			if n == 0 {
				break
			}
			if err != nil {
				panic(err)
			}
			if !passed {
				p.addBytes(buffer[0:n])
				if waitPassword(p) {
					cha <- true
					passed = true
				}
			}
			_, err = os.Stdout.Write(buffer[0:n])
			if err != nil {
				panic(err)
			}
		}
	}()

	// Start remote shell
	if err := session.Shell(); err != nil {
		panic(err)
	}
	session.Wait()
}

func waitPassword(p PatternDetector) bool {
	return bytes.Contains(p.exportBytes(), []byte("Mot de passe : ")) || bytes.Contains(p.exportBytes(), []byte("Password: "))
}
