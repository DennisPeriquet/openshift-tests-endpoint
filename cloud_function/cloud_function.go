package p

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	// Register an HTTP function with the Functions Framework
	functions.HTTP("OpenshiftTestsEndpoint", OpenshiftTestsEndpoint)
}

// OpenshiftTestsEndpoint handles the incoming HTTP request.
func OpenshiftTestsEndpoint(w http.ResponseWriter, r *http.Request) {
	auditId := sanitizeHeader(r.Header.Get("AuditId"))
	buildId := sanitizeHeader(r.Header.Get("BuildId"))

	if auditId == "" || buildId == "" {
		http.Error(w, fmt.Sprintf("Invalid request format: auditId=(%v) buildId=(%v)", auditId, buildId), http.StatusBadRequest)
		return
	}

	currentTime := time.Now().UTC()
	fmt.Printf("HTTP get received: AuditId: %s, BuildId: %s, Current Time (UTC): %s", auditId, buildId, currentTime)
	w.WriteHeader(http.StatusOK)
}

func sanitizeHeader(headerValue string) string {
	reg := regexp.MustCompile(`^[a-zA-Z0-9-_]+$`)
	if reg.MatchString(headerValue) {
		return headerValue
	}
	return ""
}
