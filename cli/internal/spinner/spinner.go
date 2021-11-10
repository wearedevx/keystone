package spinner

import (
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
)

type spin struct {
	inner *spinner.Spinner
}

type SpinnerInterface interface {
	Start() SpinnerInterface
	Stop() SpinnerInterface
}

var noSpin bool

func init() {
	v := os.Getenv("NOSPIN")

	noSpin = v == "true"
}

func buildClearString(message string) string {
	sb := new(strings.Builder)
	// Moves the cursor to the beginning of the line
	sb.WriteRune('\r')
	// Replaces the whole message with whitespaces
	for range message {
		sb.WriteRune(' ')
	}
	// Moves the cursor the beginning of the line again
	sb.WriteRune('\r')

	return sb.String()
}

// Provides a pointer to a struct implementing SpinnerInterface
func Spinner(message string) SpinnerInterface {
	s := spinner.New(
		spinner.CharSets[14],
		100*time.Millisecond,
		spinner.WithWriter(os.Stderr),
	)

	s.Suffix = " " + message
	s.FinalMSG = buildClearString(s.Suffix)

	return &spin{
		inner: s,
	}
}

// Starts the spinner
func (s *spin) Start() SpinnerInterface {
	if noSpin {
		return s
	}

	s.inner.Start()

	return s
}

// Stops the spinner
func (s *spin) Stop() SpinnerInterface {
	if noSpin {
		return s
	}

	s.inner.Stop()

	return s
}
