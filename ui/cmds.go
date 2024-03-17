package ui

import (
	"context"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

func quitProg() tea.Cmd {
	return func() tea.Msg {
		return quitProgMsg{}
	}
}

func (m model) fetchSystemVersion(system System) tea.Cmd {
	return func() tea.Msg {

		var awsConfig AWSConfig
		switch m.awsConfigSource {
		case SharedCfgProfileSrc:
			awsConfig = m.awsConfigs[getSharedProfileCfgKey(&system)]
		case DefaultCfg:
			switch system.IAMRoleToAssume {
			case "":
				awsConfig = m.awsConfigs[getDefaultCfgKey(&system)]
			default:
				awsConfig = m.awsConfigs[getRoleCfgKey(&system)]
			}
		}
		if awsConfig.err != nil {
			return processFinishedMsg{
				systemKey: system.Key,
				env:       system.Env,
				err:       awsConfig.err,
			}

		}

		ecsClient := ecs.NewFromConfig(awsConfig.config)

		services := make([]string, 1)
		services[0] = system.ServiceName
		svcs, err := ecsClient.DescribeServices(context.Background(), &ecs.DescribeServicesInput{Services: services, Cluster: &system.ClusterName})
		var version string
		if err != nil {
			return processFinishedMsg{
				systemKey: system.Key,
				env:       system.Env,
				err:       err,
			}
		}
		for _, svc := range svcs.Services {
			td := svc.TaskDefinition

			tdD, err := ecsClient.DescribeTaskDefinition(context.Background(), &ecs.DescribeTaskDefinitionInput{TaskDefinition: td})
			if err != nil {
				return processFinishedMsg{
					systemKey: system.Key,
					env:       system.Env,
					err:       err,
				}
			}
			cd := tdD.TaskDefinition.ContainerDefinitions
			for _, cdd := range cd {
				if *cdd.Name == system.ContainerName {
					versionEls := strings.Split(*cdd.Image, ":")
					if len(versionEls) > 0 {
						version = versionEls[len(versionEls)-1]
					}
					return processFinishedMsg{
						found:     true,
						systemKey: system.Key,
						env:       system.Env,
						version:   version,
					}
				}
			}
		}
		return processFinishedMsg{
			systemKey: system.Key,
			env:       system.Env,
			found:     false,
		}
	}
}
