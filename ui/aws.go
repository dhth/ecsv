package ui

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func getAWSConfigKey(system System) string {
	switch system.AWSConfigSourceType {
	case SharedCfgProfileType, AssumeRoleCfgType:
		return system.AWSConfigSource + ":" + system.AWSRegion
	default:
		return system.AWSRegion
	}
}

func getAWSConfig(system System) (aws.Config, error) {
	var cfg aws.Config
	var err error
	switch system.AWSConfigSourceType {
	case SharedCfgProfileType:
		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(system.AWSRegion),
			config.WithSharedConfigProfile(system.AWSConfigSource))
	case AssumeRoleCfgType:
		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(system.AWSRegion))
		if err != nil {
			return cfg, err
		}
		stsSvc := sts.NewFromConfig(cfg)
		creds := stscreds.NewAssumeRoleProvider(stsSvc, system.AWSConfigSource)

		cfg.Credentials = aws.NewCredentialsCache(creds)
	default:
		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(system.AWSRegion))
	}
	return cfg, err

}
