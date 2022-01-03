package loggers

import (
	"io"
	"log"
)

var loggers = makeLoggers()

func makeLoggers() []*log.Logger {
	return make([]*log.Logger, 0)
}

func AddLogger(logger *log.Logger) {
	loggers = append(loggers, logger)
}

func SetOutput(w io.Writer) {
	log.SetOutput(w)

	for _, logger := range loggers {
		logger.SetOutput(w)
	}
}
