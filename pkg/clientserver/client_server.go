package clientserver

import (
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

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
func RunServer(useHttps *bool, certFile, keyFile *string, port int) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: time.RFC3339,
		FullTimestamp:   true,
	})

	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		auditId := sanitizeHeader(r.Header.Get("Audit-ID"))

		if auditId == "" {
			http.Error(w, fmt.Sprintf("Invalid request format: Audit-ID=(%v)", auditId), http.StatusBadRequest)
			return
		}

		logger.Infof("HTTP get received: Audit-ID: %s", auditId)
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
func RunClient(count, port int) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: time.RFC3339,
		FullTimestamp:   true,
	})

	for i := 1; i <= count; i++ {
		go func(clientID int) {
			ticker := time.NewTicker(1 * time.Second)
			for range ticker.C {
				sendRequest(clientID, logger, port)
			}
		}(i)
	}

	select {}
}

// sendRequest sends an HTTP GET request to the server and logs the response.
// This is just a way to simulate client requests to the server.
func sendRequest(clientID int, logger *logrus.Logger, port int) {
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
	req.Header.Add("Audit-ID", strconv.Itoa(randomNumber))

	resp, err := client.Do(req)
	if err != nil {
		logger.Error(err)
		return
	}
	defer resp.Body.Close()

	logger.Infof("Client %d received response: %s", clientID, resp.Status)
}
