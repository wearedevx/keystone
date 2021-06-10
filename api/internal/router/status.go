package router

import "net/http"

type doneWriter struct {
	http.ResponseWriter
	done bool
}

func (w *doneWriter) WriteHeader(status int) {
	if !w.done {
		w.done = true
		w.ResponseWriter.WriteHeader(status)
	}
}

func (w *doneWriter) Write(b []byte) (int, error) {
	w.done = true
	return w.ResponseWriter.Write(b)
}

func newDoneWriter(w http.ResponseWriter) *doneWriter {
	return &doneWriter{ResponseWriter: w, done: false}
}
