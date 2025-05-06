package ui

import "github.com/dhth/ecsv/internal/types"

type HTMLDataRow struct {
	Data   []string
	InSync bool
}
type HTMLData struct {
	Favicon   string
	Title     string
	TitleURL  string
	Columns   []string
	Rows      []HTMLDataRow
	Errors    []error
	Timestamp string
}

type Config struct {
	EnvSequence         []string
	SystemKeys          []string
	OutputFmt           types.OutputFmt
	HTMLTemplate        string
	HTMLTitle           string
	HTMLTitleURL        string
	HTMLFaviconContents []byte
	Style               types.TableStyle
	ShowRegisteredAt    bool
}

type SystemResult struct {
	SystemKey string
	Env       string
	Version   string
	Found     bool
	Err       error
}

type ConfigSourceType uint

const (
	DefaultCfgType ConfigSourceType = iota
	SharedCfgProfileType
	AssumeRoleCfgType
)
