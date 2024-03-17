package ui

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func getSharedProfileCfgKey(system *System) string {
	return system.AWSProfile + ":" + system.AWSRegion
}

func getRoleCfgKey(system *System) string {
	return system.IAMRoleToAssume + ":" + system.AWSRegion
}

func getDefaultCfgKey(system *System) string {
	return system.AWSRegion
}

func getAWSConfig(profile string, region string) (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithSharedConfigProfile(profile))
	return cfg, err

}

func getDefaultConfig(region string) (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region))
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}

func getRoleConfig(roleArn string, region string) (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region))
	if err != nil {
		return cfg, err
	}

	stsSvc := sts.NewFromConfig(cfg)
	creds := stscreds.NewAssumeRoleProvider(stsSvc, roleArn)

	cfg.Credentials = aws.NewCredentialsCache(creds)
	return cfg, nil
}
