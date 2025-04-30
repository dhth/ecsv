package types

import (
	"sort"
	"time"
)

type OutputFmt uint

const (
	DefaultFmt OutputFmt = iota
	TabularFmt
	HTMLFmt
)

func OutputFormats() []string {
	return []string{"default", "table", "html"}
}

type AWSConfigSourceType uint

const (
	DefaultCfgType AWSConfigSourceType = iota
	SharedCfgProfileType
	AssumeRoleCfgType
)

type System struct {
	Key                 string
	Env                 string
	AWSConfigSourceType AWSConfigSourceType
	AWSConfigSource     string
	AWSRegion           string
	IAMRoleToAssume     string
	ClusterName         string
	ServiceName         string
	ContainerName       string
}

func (s System) AWSConfigKey() string {
	switch s.AWSConfigSourceType {
	case SharedCfgProfileType, AssumeRoleCfgType:
		return s.AWSConfigSource + ":" + s.AWSRegion
	default:
		return s.AWSRegion
	}
}

type SystemResult struct {
	SystemKey    string
	Env          string
	Version      string
	Found        bool
	RegisteredAt *time.Time
	Err          error
}

type TableStyle string

func (ts TableStyle) String() string {
	return string(ts)
}

const (
	ASCIIStyle TableStyle = "ascii"
	BlankStyle TableStyle = "blank"
	DotsStyle  TableStyle = "dots"
	SharpStyle TableStyle = "sharp"
)

var styles = map[string]TableStyle{
	string(ASCIIStyle): ASCIIStyle,
	string(BlankStyle): BlankStyle,
	string(DotsStyle):  DotsStyle,
	string(SharpStyle): SharpStyle,
}

func TableStyleStrings() []string {
	values := make([]string, 0, len(styles))
	for k := range styles {
		values = append(values, k)
	}
	sort.Strings(values)
	return values
}

func GetStyle(styleStr string) (TableStyle, bool) {
	style, ok := styles[styleStr]
	return style, ok
}
