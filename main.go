package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
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
		runServer(useHTTPS, certFile, keyFile)
	} else if *mode == "client" {
		runClient(*clientCount)
	} else {
		fmt.Printf("Usage: %s -mode [client|server] [-count n] [-https]\n", os.Args[0])
	}
}

// sanitizeHeader removes any characters from the header value that are not
// alphanumeric, hyphens, or underscores.
func sanitizeHeader(headerValue string) string {
	reg := regexp.MustCompile(`^[a-zA-Z0-9-_]+$`)
	if reg.MatchString(headerValue) {
		return headerValue
	}
	return ""
}

// runServer is the listens to HTTP requests and logs the request headers.
func runServer(useHttps *bool, certFile, keyFile *string) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: time.RFC3339,
		FullTimestamp:   true,
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		auditId := sanitizeHeader(r.Header.Get("AuditId"))
		buildId := sanitizeHeader(r.Header.Get("BuildId"))

		if auditId == "" || buildId == "" {
			http.Error(w, fmt.Sprintf("Invalid request format: auditId=(%v) buildId=(%v)", auditId, buildId), http.StatusBadRequest)
			return
		}

		logger.Infof("HTTP get received: AuditId: %s, BuildId: %s", auditId, buildId)
		w.WriteHeader(http.StatusOK)
	})

	if *useHttps {
		logger.Infof("Starting HTTPS server on port %d", port)
		err := http.ListenAndServeTLS(fmt.Sprintf(":%d", port), *certFile, *keyFile, nil)
		if err != nil {
			logger.Fatalf("Failed to start HTTPS server: %v", err)
		}
	} else {
		logger.Infof("Starting HTTP server on port %d", port)
		err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
		if err != nil {
			logger.Fatalf("Failed to start HTTP server: %v", err)
		}
	}
}

// runClient runs a number of clients (n) that send HTTP requests to the server.
// This is just a test client to help debug the server.
func runClient(count int) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: time.RFC3339,
		FullTimestamp:   true,
	})

	for i := 1; i <= count; i++ {
		go func(clientID int) {
			ticker := time.NewTicker(1 * time.Second)
			for range ticker.C {
				sendRequest(clientID, logger)
			}
		}(i)
	}

	select {}
}

// sendRequest sends an HTTP GET request to the server and logs the response.
// This is just a way to simulate client requests to the server.
func sendRequest(clientID int, logger *logrus.Logger) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d/health", port), nil)
	if err != nil {
		logger.Error(err)
		return
	}

	// generate a randome number from 1 to 1000000
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	randomNumber := random.Intn(1000000) + 1
	req.Header.Add("AuditId", strconv.Itoa(randomNumber))

	randomNumber = random.Intn(3) + 2
	req.Header.Add("BuildId", "build0"+strconv.Itoa(randomNumber))

	resp, err := client.Do(req)
	if err != nil {
		logger.Error(err)
		return
	}
	defer resp.Body.Close()

	logger.Infof("Client %d received response: %s", clientID, resp.Status)
}
