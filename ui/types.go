package ui

import (
	"github.com/aws/aws-sdk-go-v2/aws"
)

type AWSConfig struct {
	config aws.Config
	err    error
}

type System struct {
	Key           string
	Env           string
	AWSProfile    string
	AWSRegion     string
	ClusterName   string
	ServiceName   string
	ContainerName string
}
