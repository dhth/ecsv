package ui

import (
	"github.com/aws/aws-sdk-go-v2/aws"
)

type AWSConfig struct {
	config aws.Config
	err    error
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

type OutFormat uint

const (
	UnspecifiedFmt OutFormat = iota
	PlainTextFmt
	HTMLFmt
)

type HTMLDataRow struct {
	Data   []string
	InSync bool
}
type HTMLData struct {
	Title     string
	Columns   []string
	Rows      []HTMLDataRow
	Errors    *[]error
	Timestamp string
}
