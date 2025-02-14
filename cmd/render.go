package cmd

import (
	"fmt"

	"github.com/dhth/ecsv/internal/aws"
	"github.com/dhth/ecsv/internal/types"
	"github.com/dhth/ecsv/internal/ui"
)

func render(systems []types.System, config ui.Config, awsConfigs map[string]aws.Config) error {
	results := make(map[string]map[string]types.SystemResult)
	resultChannel := make(chan types.SystemResult)

	counter := 0
	for _, s := range systems {
		awsConfig := awsConfigs[s.AWSConfigKey()]
		if results[s.Key] == nil {
			results[s.Key] = make(map[string]types.SystemResult)
		}
		results[s.Key][s.Env] = types.SystemResult{}

		if awsConfig.Err != nil {
			results[s.Key][s.Env] = types.SystemResult{
				SystemKey: s.Key,
				Env:       s.Env,
				Err:       awsConfig.Err,
			}
			continue
		}
		go func(system types.System) {
			resultChannel <- aws.FetchSystemVersion(system, awsConfig)
		}(s)
		counter++
	}

	for range counter {
		r := <-resultChannel
		results[r.SystemKey][r.Env] = r
	}

	output, err := ui.GetOutput(config, results)
	if err != nil {
		return err
	}
	fmt.Print(output)
	return nil
}
