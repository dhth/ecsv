package aws

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/dhth/ecsv/internal/types"
)

type Config struct {
	Config aws.Config
	Err    error
}

func GetConfig(system types.VersionsConfig) (aws.Config, error) {
	var cfg aws.Config
	var err error
	switch system.AWSConfigSourceType {
	case types.SharedCfgProfileType:
		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(system.AWSRegion),
			config.WithSharedConfigProfile(system.AWSConfigSource))
	case types.AssumeRoleCfgType:
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

func FetchSystemVersion(system types.VersionsConfig, awsConfig Config) types.VersionResult {
	ecsClient := ecs.NewFromConfig(awsConfig.Config)

	services := make([]string, 1)
	services[0] = system.ServiceName
	svcs, err := ecsClient.DescribeServices(context.Background(), &ecs.DescribeServicesInput{Services: services, Cluster: &system.ClusterName})
	if err != nil {
		return types.VersionResult{
			SystemKey: system.Key,
			Env:       system.Env,
			Err:       err,
		}
	}
	for _, svc := range svcs.Services {
		td := svc.TaskDefinition

		describeTDOutput, err := ecsClient.DescribeTaskDefinition(context.Background(), &ecs.DescribeTaskDefinitionInput{TaskDefinition: td})
		if err != nil {
			return types.VersionResult{
				SystemKey: system.Key,
				Env:       system.Env,
				Err:       err,
			}
		}
		containerDefs := describeTDOutput.TaskDefinition.ContainerDefinitions
		for _, containerDef := range containerDefs {
			if *containerDef.Name == system.ContainerName {
				versionEls := strings.Split(*containerDef.Image, ":")
				var version string
				if len(versionEls) > 0 {
					version = versionEls[len(versionEls)-1]
				}
				var registeredAt *time.Time
				if describeTDOutput != nil && describeTDOutput.TaskDefinition != nil {
					registeredAt = describeTDOutput.TaskDefinition.RegisteredAt
				}

				return types.VersionResult{
					Found:        true,
					SystemKey:    system.Key,
					Env:          system.Env,
					Version:      version,
					RegisteredAt: registeredAt,
				}
			}
		}
	}
	return types.VersionResult{
		SystemKey: system.Key,
		Env:       system.Env,
		Found:     false,
	}
}
