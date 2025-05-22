package ui

import (
	"fmt"
	"strings"

	"github.com/dhth/ecsv/internal/types"
)

type VersionRow struct {
	Data   []string
	InSync bool
}

type HTMLData struct {
	Title     string
	TitleURL  string
	Columns   []string
	Rows      []VersionRow
	Changes   []types.ChangesResult
	Errors    []error
	Timestamp string
}

type Config struct {
	EnvSequence      []string
	SystemKeys       []string
	OutputFmt        types.OutputFmt
	HTMLConfig       HTMLOutputConfig
	TableConfig      TableOutputConfig
	ShowRegisteredAt bool
}

type HTMLOutputConfig struct {
	Template string
	Title    string
	TitleURL string
	Open     bool
}

type TableOutputConfig struct {
	Style types.TableStyle
}

func (c Config) String() string {
	switch c.OutputFmt {
	case types.HTMLFmt:
		return strings.TrimSpace(fmt.Sprintf(`
- env sequence          %v
- system keys           %v
- output format         %s
- html title            %s
- html title url        %s
- show registererd url  %v
`,
			c.EnvSequence,
			c.SystemKeys,
			c.OutputFmt.String(),
			c.HTMLConfig.Title,
			c.HTMLConfig.TitleURL,
			c.ShowRegisteredAt,
		))
	case types.TabularFmt:
		return strings.TrimSpace(fmt.Sprintf(`
- env sequence          %v
- system keys           %v
- output format         %s
- style                 %s
`,
			c.EnvSequence,
			c.SystemKeys,
			c.OutputFmt.String(),
			c.TableConfig.Style.String(),
		))
	default:
		return strings.TrimSpace(fmt.Sprintf(`
- env sequence          %v
- system keys           %v
- output format         %s
- show registererd url  %v
`,
			c.EnvSequence,
			c.SystemKeys,
			c.OutputFmt.String(),
			c.ShowRegisteredAt,
		))
	}
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
