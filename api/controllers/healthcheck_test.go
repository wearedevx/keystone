package controllers

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/julienschmidt/httprouter"
)

type mockResponseWriter struct {
	status  int
	headers map[string][]string
	body    *bytes.Buffer
}

func newMockResponse() *mockResponseWriter {
	return &mockResponseWriter{
		status:  0,
		headers: make(map[string][]string),
		body:    new(bytes.Buffer),
	}
}

// Header returns the header map that will be sent by
// WriteHeader.
func (m *mockResponseWriter) Header() http.Header {
	return m.headers
}

// Write writes the data to the connection as part of an HTTP reply.
func (m *mockResponseWriter) Write(content []byte) (int, error) {
	return m.body.Write(content)
}

// WriteHeader sends an HTTP response header with the provided
func (m *mockResponseWriter) WriteHeader(statusCode int) {
	m.status = statusCode
}

func TestGetHealthCheck(t *testing.T) {
	type args struct {
		w   http.ResponseWriter
		in1 *http.Request
		in2 httprouter.Params
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
	}{
		{
			name: "it return ok",
			args: args{
				w:   newMockResponse(),
				in1: &http.Request{},
				in2: []httprouter.Param{},
			},
			wantStatus: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			GetHealthCheck(tt.args.w, tt.args.in1, tt.args.in2)

			//got := tt.args.w.(*mockResponseWriter)
			//if got.status != tt.wantStatus {
			// TODO
			//}
		})
	}
}
