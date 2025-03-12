package cmd

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/dhth/ecsv/internal/types"

	"gopkg.in/yaml.v3"
)

var (
	errInvalidConfigSourceProvided = errors.New("invalid aws-system-source provided")
	errConfigIsInvalidYAML         = errors.New("config file is not valid YAML")
)

type Config struct {
	EnvSequence []string `yaml:"env-sequence"`
	Systems     []struct {
		Key  string `yaml:"key"`
		Envs []struct {
			Name            string `yaml:"name"`
			AwsConfigSource string `yaml:"aws-config-source"`
			AwsRegion       string `yaml:"aws-region"`
			Cluster         string `yaml:"cluster"`
			Service         string `yaml:"service"`
			ContainerName   string `yaml:"container-name"`
		} `yaml:"envs"`
	} `yaml:"systems"`
}

func expandTilde(path string, homeDir string) string {
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(homeDir, path[2:])
	}
	return path
}

func readConfig(configBytes []byte, keyRegex *regexp.Regexp) ([]string, []types.System, error) {
	cfg := Config{}
	err := yaml.Unmarshal(configBytes, &cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %s", errConfigIsInvalidYAML, err.Error())
	}

	var systems []types.System

	for _, system := range cfg.Systems {
		if keyRegex != nil && !keyRegex.Match([]byte(system.Key)) {
			continue
		}

		for _, env := range system.Envs {

			var awsConfigType types.AWSConfigSourceType
			var awsConfigSource string
			switch {
			case env.AwsConfigSource == "default":
				awsConfigType = types.DefaultCfgType
			case strings.HasPrefix(env.AwsConfigSource, "profile:::"):
				configElements := strings.Split(env.AwsConfigSource, "profile:::")
				awsConfigSource = configElements[len(configElements)-1]
				awsConfigType = types.SharedCfgProfileType
			case strings.HasPrefix(env.AwsConfigSource, "assume-role:::"):
				configElements := strings.Split(env.AwsConfigSource, "assume-role:::")
				awsConfigSource = configElements[len(configElements)-1]
				awsConfigType = types.AssumeRoleCfgType
			default:
				return nil,
					nil,
					fmt.Errorf("%w: system: %s env: %s", errInvalidConfigSourceProvided,
						system.Key,
						env.Name)
			}
			systems = append(systems, types.System{
				Key:                 system.Key,
				Env:                 env.Name,
				AWSConfigSourceType: awsConfigType,
				AWSConfigSource:     awsConfigSource,
				AWSRegion:           env.AwsRegion,
				ClusterName:         env.Cluster,
				ServiceName:         env.Service,
				ContainerName:       env.ContainerName,
			})
		}
	}
	return cfg.EnvSequence, systems, nil
}
