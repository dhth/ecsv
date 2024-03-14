package ui

type processFinishedMsg struct {
	systemKey string
	env       string
	version   string
	found     bool
	err       error
}

type quitProgMsg struct{}
