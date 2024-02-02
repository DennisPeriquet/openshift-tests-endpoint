package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/DennisPeriquet/openshift-tests-endpoint/pkg/clientserver"
)

const (
	port           = 49888
	defaultCertPem = "./cert.pem"
	defaultKeyPem  = "./key.pem"
)

func main() {
	mode := flag.String("mode", "", "run in 'client' or 'server' mode")
	clientCount := flag.Int("count", 1, "number of clients to run in client mode")
	useHTTPS := flag.Bool("https", false, "use HTTPS (default HTTP)")
	certFile := flag.String("cert", defaultCertPem, "TLS certificate file")
	keyFile := flag.String("key", defaultKeyPem, "TLS private key file")

	flag.Parse()

	if *mode == "server" {
		clientserver.RunServer(useHTTPS, certFile, keyFile, port)
	} else if *mode == "client" {
		clientserver.RunClient(*clientCount, port)
	} else {
		fmt.Printf("Usage: %s -mode [client|server] [-count n] [-https]\n", os.Args[0])
	}
}
