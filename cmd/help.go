package cmd

import "fmt"

var (
	configSampleFormat = `
env-sequence: ["qa", "staging"]
systems:
- key: service-a
  envs:
  - name: qa
    aws-profile: qa
    aws-region: eu-central-1
    cluster: 1brd-qa
    service: service-a-fargate
    container-name: service-a-qa-Service
  - name: staging
    aws-profile: qa
    aws-region: eu-central-1
    cluster: 1brd-staging
    service: service-a-fargate
    container-name: service-a-staging-Service
- key: service-b
  envs:
  - name: qa
    aws-profile: qa
    aws-region: eu-central-1
    cluster: 1brd-qa
    service: service-b-fargate
    container-name: service-b-qa-Service
  - name: staging
    aws-profile: qa
    aws-region: eu-central-1
    cluster: 1brd-staging
    service: service-b-fargate
    container-name: service-b-staging-Service
`
	helpText = `Quickly check the code versions of containers running in your ECS services across various environments.

Usage: ecsv [flags]`
)

func cfgErrSuggestion(msg string) string {
	return fmt.Sprintf(`%s

Make sure to structure the yml config file as follows:

%s

Use "ecsv -help" for more information`,
		msg,
		configSampleFormat,
	)
}
