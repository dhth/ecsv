package cmd

import (
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
			Name          string `yaml:"name"`
			AwsProfile    string `yaml:"aws-profile"`
			AwsRegion     string `yaml:"aws-region"`
			Cluster       string `yaml:"cluster"`
			Service       string `yaml:"service"`
			ContainerName string `yaml:"container-name"`
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
			systems = append(systems, ui.System{
				Key:           system.Key,
				Env:           env.Name,
				AWSProfile:    env.AwsProfile,
				AWSRegion:     env.AwsRegion,
				ClusterName:   env.Cluster,
				ServiceName:   env.Service,
				ContainerName: env.ContainerName,
			})
		}
	}
	return t.EnvSequence, systems, err

}
