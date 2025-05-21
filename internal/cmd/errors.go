package cmd

import (
	"errors"

	"github.com/dhth/ecsv/internal/ui"
)

type ErrorFollowUp struct {
	IsUnexpected bool
	Message      string
}

func GetErrorFollowUp(err error) (ErrorFollowUp, bool) {
	var zero ErrorFollowUp

	if errors.Is(err, ui.ErrCouldntCreateTable) || errors.Is(err, ui.ErrCouldntParseBuiltInHTMLTemplate) {
		return unexpectedErr("")
	} else if errors.Is(err, ui.ErrCouldntParseHTMLTemplate) {
		return expectedErr("Maybe take a look at ecsv's built in template (on GitHub)")
	}

	return zero, false
}

func unexpectedErr(message string) (ErrorFollowUp, bool) {
	return ErrorFollowUp{
		IsUnexpected: true,
		Message:      message,
	}, true
}

func expectedErr(message string) (ErrorFollowUp, bool) {
	return ErrorFollowUp{
		IsUnexpected: false,
		Message:      message,
	}, true
}
