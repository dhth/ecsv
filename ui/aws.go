package ui

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func getAWSCfgKey(system *System) string {
	return system.AWSProfile + ":" + system.AWSRegion
}

func getAWSConfig(profile string, region string) (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithSharedConfigProfile(profile))
	return cfg, err

}
