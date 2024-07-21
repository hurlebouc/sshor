package ssh

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/hurlebouc/sshor/config"
	"github.com/things-go/go-socks5"
)

func Socks(hostConf config.Host, options Options, passwordFlag, keepassPwdFlag string, listeningIp string, listeningPort uint16) {
	sshClient := getSshClient(hostConf, passwordFlag, keepassPwdFlag)
	defer sshClient.Close()
	// Create a SOCKS5 server
	logger := log.New(os.Stdout, "socks5: ", log.LstdFlags)
	server := socks5.NewServer(
		socks5.WithLogger(socks5.NewLogger(logger)),
		socks5.WithDial(func(ctx context.Context, network, addr string) (net.Conn, error) {
			if options.Verbose {
				logger.Printf("(%s) addr: %s", network, addr)
			}
			return sshClient.Dial(network, addr)
		}),
	)

	listener := listen(localNet{}, listeningIp, listeningPort)

	// Create SOCKS5 proxy on localhost port 8000
	if err := server.Serve(listener); err != nil {
		panic(err)
	}
}
