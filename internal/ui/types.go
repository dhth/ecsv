package ui

import (
	"fmt"
	"strings"

	"github.com/dhth/ecsv/internal/types"
)

type HTMLDataRow struct {
	Data   []string
	InSync bool
}
type HTMLData struct {
	Title     string
	TitleURL  string
	Columns   []string
	Rows      []HTMLDataRow
	Errors    []error
	Timestamp string
}

type Config struct {
	EnvSequence      []string
	SystemKeys       []string
	OutputFmt        types.OutputFmt
	HTMLTemplate     string
	HTMLTitle        string
	HTMLTitleURL     string
	Style            types.TableStyle
	ShowRegisteredAt bool
}

func (c Config) String() string {
	return strings.TrimSpace(fmt.Sprintf(`
- env sequence          %v
- system keys           %v
- output format         %s
- html title            %s
- html title url        %s
- style                 %s
- show registererd url  %v
`,
		c.EnvSequence,
		c.SystemKeys,
		c.OutputFmt.String(),
		c.HTMLTitle,
		c.HTMLTitleURL,
		c.Style.String(),
		c.ShowRegisteredAt,
	))
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
