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

func ForwardProxy(hostConf config.Host, options Options, passwordFlag, keepassPwdFlag string, listeningIp string, listeningPort uint16, destinationAddr string, destinationPort uint16) {
	keepassPwdMap := InitKeepassPwdMap(hostConf, keepassPwdFlag)
	ctx := InitContext()
	sshClient, _ := NewSshClient(ctx, hostConf, passwordFlag, keepassPwdMap)
	defer sshClient.Close()
	if sshClient.Client == nil {
		log.Panicln("Cannot change user of proxied connection")
	}

	listPort := listeningPort
	if listeningPort == 0 {
		listPort = randomPort()
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", listeningIp, listPort))

	if listeningPort == 0 {
		for err != nil {
			listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", listeningIp, listeningPort))
		}
	}

	log.Printf("listening at %s:%d", listeningIp, listPort)

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

			remoteConn, err := sshClient.Client.Dial("tcp", fmt.Sprintf("%s:%d", destinationAddr, destinationPort))
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
