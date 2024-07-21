package ssh

import (
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"net"

	"github.com/hurlebouc/sshor/config"
)

func randomPort() uint16 {
	return uint16(rand.Uint32())
}

type network interface {
	Listen(network string, address string) (net.Listener, error)
	Dial(n string, addr string) (net.Conn, error)
}

type localNet struct{}

func (ln localNet) Listen(network string, address string) (net.Listener, error) {
	return net.Listen(network, address)
}

func (ln localNet) Dial(n string, addr string) (net.Conn, error) {
	return net.Dial(n, addr)
}

func listen(srcNet network, listeningIp string, listeningPort uint16) net.Listener {
	listPort := listeningPort
	if listeningPort == 0 {
		listPort = randomPort()
	}

	listener, err := srcNet.Listen("tcp", fmt.Sprintf("%s:%d", listeningIp, listPort))

	if listeningPort == 0 {
		for err != nil {
			listPort = randomPort()
			listener, err = srcNet.Listen("tcp", fmt.Sprintf("%s:%d", listeningIp, listPort))
		}
	}

	log.Printf("listening at %s:%d", listeningIp, listPort)
	return listener
}

func proxy(options Options, srcNet, dstNet network, listeningIp string, listeningPort uint16, destinationAddr string, destinationPort uint16) {

	listener := listen(srcNet, listeningIp, listeningPort)

	for {
		localConn, err := listener.Accept()
		if err != nil {
			log.Panic(err)
		}
		if options.Verbose {
			log.Printf("new connection from %s", localConn.LocalAddr().String())
		}

		go func() {
			defer localConn.Close()

			remoteConn, err := dstNet.Dial("tcp", fmt.Sprintf("%s:%d", destinationAddr, destinationPort))
			if err != nil {
				panic(err)
			}
			defer remoteConn.Close()
			forwardChan := make(chan bool)
			go func() {
				io.Copy(remoteConn, localConn)
				close(forwardChan)
			}()
			backwardChan := make(chan bool)
			go func() {
				io.Copy(localConn, remoteConn)
				close(backwardChan)
			}()
			for range forwardChan {
			}
			for range backwardChan {
			}
			if options.Verbose {
				log.Printf("close connection from %s", localConn.LocalAddr().String())
			}
		}()

	}
}

func ForwardProxy(hostConf config.Host, options Options, passwordFlag, keepassPwdFlag string, listeningIp string, listeningPort uint16, destinationAddr string, destinationPort uint16) {
	sshClient := getSshClient(hostConf, passwordFlag, keepassPwdFlag)
	defer sshClient.Close()
	proxy(options, localNet{}, sshClient, listeningIp, listeningPort, destinationAddr, destinationPort)
}
func BackwardProxy(hostConf config.Host, options Options, passwordFlag, keepassPwdFlag string, listeningIp string, listeningPort uint16, destinationAddr string, destinationPort uint16) {
	sshClient := getSshClient(hostConf, passwordFlag, keepassPwdFlag)
	defer sshClient.Close()
	proxy(options, sshClient, localNet{}, listeningIp, listeningPort, destinationAddr, destinationPort)
}
