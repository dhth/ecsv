package types

type OutputFmt uint

const (
	DefaultFmt OutputFmt = iota
	TabularFmt
	HTMLFmt
)

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
	SystemKey string
	Env       string
	Version   string
	Found     bool
	Err       error
}
