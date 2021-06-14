package ui

import (
	"time"

	"github.com/briandowns/spinner"
)

type Spin struct {
	Spinner *spinner.Spinner
}

var Spinner = &Spin{
	spinner.New(spinner.CharSets[14], 80*time.Millisecond),
}

func GetSpinner() *Spin {
	return &Spin{
		spinner.New(spinner.CharSets[14], 80*time.Millisecond),
	}
}

func (s *Spin) Start(prefix string, suffix string) {
	s.UpdateText(prefix, suffix)
	s.Spinner.Start()
}

func (s *Spin) StartEmpty() {
	s.Spinner.Start()
}

func (s *Spin) Stop() {
	s.Spinner.Stop()
	s.UpdateText("", "")
}

func (s *Spin) Pause() {
	s.Spinner.Stop()
}

func (s *Spin) UpdateText(prefix string, suffix string) {
	s.Spinner.Prefix = prefix
	s.Spinner.Suffix = suffix
}
