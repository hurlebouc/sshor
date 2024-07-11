package cmd

var proxyOptions struct {
	listeningIP     string
	listeningPort   uint16
	destinationIP   string
	destinationPort uint16
}

const LISTENING_IP = "listening-ip"
const LISTENING_PORT = "listening-port"
const DESTINATION_ADDR = "destination-addr"
const DESTINATION_PORT = "destination-port"
