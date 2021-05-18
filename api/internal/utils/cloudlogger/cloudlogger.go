package cloudlogger

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type LogEntry struct {
	Message  string `json:"message"`
	Severity string `json:"severity,omitempty"`
	Trace    string `json:"logging.googleapis.com/trace,omitempty"`
}

func writeLog(r *http.Request, severity string, message string, args ...interface{}) {
	var trace string
	traceHeader := r.Header.Get("X-Cloud-Trace-Context")
	traceParts := strings.Split(traceHeader, "/")
	if len(traceParts) > 0 && len(traceParts[0]) > 0 {
		trace = fmt.Sprintf("projects/%s/traces/%s", "wearedevx", traceParts[0])
	}

	log.Println(LogEntry{
		Message:  fmt.Sprintf(message, args...),
		Severity: severity,
		Trace:    trace,
	}.String())
}

func Init() {
	log.SetFlags(0)
	log.SetPrefix("")
}

func Notice(r *http.Request, message string, args ...interface{}) {
	writeLog(r, "NOTICE", message, args...)
}

func Info(r *http.Request, message string, args ...interface{}) {
	writeLog(r, "INFO", message, args...)
}

func Error(r *http.Request, message string, args ...interface{}) {
	writeLog(r, "ERROR", message, args...)
}

// String renders an entry structure to the JSON format expected by Stackdriver.
func (e LogEntry) String() string {
	if e.Severity == "" {
		e.Severity = "INFO"
	}
	out, err := json.Marshal(e)
	if err != nil {
	}
	return string(out)
}
