package cmd

import (
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/dhth/ecsv/ui"

	"gopkg.in/yaml.v3"
)

type T struct {
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

func expandTilde(path string) string {
	if strings.HasPrefix(path, "~") {
		usr, err := user.Current()
		if err != nil {
			return path
		}
		return strings.Replace(path, "~", usr.HomeDir, 1)
	}
	return path
}

func readConfig(filePath string) ([]string, []ui.System, error) {
	localFile, err := os.ReadFile(filePath)
	if err != nil {
		os.Exit(1)
	}
	t := T{}
	err = yaml.Unmarshal(localFile, &t)
	if err != nil {
		return nil, nil, err
	}

	systems := make([]ui.System, 0)

	for _, system := range t.Systems {
		for _, env := range system.Envs {

			var awsConfigType ui.AWSConfigSourceType
			var awsConfigSource string
			switch {
			case env.AwsConfigSource == "default":
				awsConfigType = ui.DefaultCfgType
			case strings.HasPrefix(env.AwsConfigSource, "profile:::"):
				configElements := strings.Split(env.AwsConfigSource, "profile:::")
				awsConfigSource = configElements[len(configElements)-1]
				awsConfigType = ui.SharedCfgProfileType
			case strings.HasPrefix(env.AwsConfigSource, "assume-role:::"):
				configElements := strings.Split(env.AwsConfigSource, "assume-role:::")
				awsConfigSource = configElements[len(configElements)-1]
				awsConfigType = ui.AssumeRoleCfgType
			default:
				return nil,
					nil,
					fmt.Errorf("system with key %s doesn't have a valid aws-config-source for env %s",
						system.Key,
						env.Name)
			}
			systems = append(systems, ui.System{
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
	return t.EnvSequence, systems, err

}
