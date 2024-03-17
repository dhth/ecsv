package cmd

import (
	"errors"
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
			Name            string  `yaml:"name"`
			AwsProfile      *string `yaml:"aws-profile"`
			AwsRegion       string  `yaml:"aws-region"`
			IAMRoleToAssume *string `yaml:"iam-role-to-assume"`
			Cluster         string  `yaml:"cluster"`
			Service         string  `yaml:"service"`
			ContainerName   string  `yaml:"container-name"`
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

func readConfig(filePath string, awsConfigSource ui.AWSConfigSource) ([]string, []ui.System, error) {
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
			if awsConfigSource == ui.SharedCfgProfileSrc {
				if env.AwsProfile == nil {
					return nil,
						nil,
						errors.New(fmt.Sprintf("system with key %s doesn't have an AWS profile set for env %s, which is needed when when using shared AWS profile configuration",
							system.Key,
							env.Name))
				}
			}
			var awsProfile string
			if awsConfigSource == ui.SharedCfgProfileSrc {
				awsProfile = *env.AwsProfile
			}

			var iamRoleToAssume string
			if env.IAMRoleToAssume != nil {
				iamRoleToAssume = *env.IAMRoleToAssume
			}
			systems = append(systems, ui.System{
				Key:             system.Key,
				Env:             env.Name,
				AWSProfile:      awsProfile,
				AWSRegion:       env.AwsRegion,
				IAMRoleToAssume: iamRoleToAssume,
				ClusterName:     env.Cluster,
				ServiceName:     env.Service,
				ContainerName:   env.ContainerName,
			})
		}
	}
	return t.EnvSequence, systems, err

}
