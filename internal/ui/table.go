package ui

import "github.com/dhth/ecsv/internal/types"

type tableSeparators struct {
	center string
	column string
	row    string
}

func getTableSeparators(style types.TableStyle) tableSeparators {
	switch style {
	case types.ASCIIStyle:
		return tableSeparators{"+", "|", "-"}
	case types.BlankStyle:
		return tableSeparators{" ", " ", " "}
	case types.DotsStyle:
		return tableSeparators{":", ":", "·"}
	default:
		return tableSeparators{"┼", "│", "─"}
	}
}
