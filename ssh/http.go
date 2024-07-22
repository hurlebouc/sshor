package ssh

import (
	"net"
	"net/http"

	"github.com/elazarl/goproxy"
	"github.com/hurlebouc/sshor/config"
)

func Http(hostConf config.Host, options Options, passwordFlag, keepassPwdFlag string, listeningIp string, listeningPort uint16) {
	sshClient := getSshClient(hostConf, passwordFlag, keepassPwdFlag)
	defer sshClient.Close()
	// Create a HTTP Proxy server
	proxy := goproxy.NewProxyHttpServer()
	proxy.ConnectDial = func(network, addr string) (net.Conn, error) {
		return sshClient.Dial(network, addr)
	}
	proxy.Verbose = options.Verbose

	listener := listen(localNet{}, listeningIp, listeningPort)

	// Create SOCKS5 proxy on localhost port 8000
	if err := http.Serve(listener, proxy); err != nil {
		panic(err)
	}
}
