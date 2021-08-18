package spinner

import (
	"os"
	"time"

	"github.com/briandowns/spinner"
)

type spin struct {
	inner *spinner.Spinner
}

type SpinnerInterface interface {
	Start()
	Stop()
}

var noSpin bool

func init() {
	v := os.Getenv("NOSPIN")

	noSpin = v == "true"
}

// Provides a pointer to a struct implementing SpinnerInterface
func Spinner(message string) SpinnerInterface {
	s := spinner.New(
		spinner.CharSets[14],
		100*time.Millisecond,
		spinner.WithWriter(os.Stderr),
	)

	s.Suffix = message
	s.FinalMSG = "\r"

	return &spin{
		inner: s,
	}
}

// Starts the spinner
func (s *spin) Start() {
	if noSpin {
		return
	}

	s.inner.Start()
}

// Stops the spinner
func (s *spin) Stop() {
	if noSpin {
		return
	}

	s.inner.Stop()
}
