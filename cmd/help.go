package cmd

import (
	"fmt"
)

var (
	configSampleFormat = `
env-sequence: ["qa", "staging"]
systems:
- key: service-a
  envs:
  - name: qa
    aws-config-source: profile:::qa
    aws-region: eu-central-1
    cluster: 1brd-qa
    service: service-a-fargate
    container-name: service-a-qa-Service
  - name: staging
    aws-config-source: profile:::staging
    aws-region: eu-central-1
    cluster: 1brd-staging
    service: service-a-fargate
    container-name: service-a-staging-Service
- key: service-b
  envs:
  - name: qa
    aws-config-source: assume-role:::arn:aws:iam::XXX:role/your-role-name
    aws-region: eu-central-1
    cluster: 1brd-qa
    service: service-b-fargate
    container-name: service-b-qa-Service
  - name: staging
    aws-config-source: default
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

Make sure to structure the yml config as follows:

%s

Config source (aws-config-source):

- profile:::qa
    will fetch AWS config from the local shared profile specified.
- assume-role:::arn:aws:iam::XXX:role/your-role-name
    will assume the role specified and use that for fetching AWS config
- default
    will use the default local AWS config

Use "ecsv -help" for more information`,
		msg,
		configSampleFormat,
	)
}
