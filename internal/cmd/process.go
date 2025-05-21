package cmd

import (
	"fmt"
	"sync"

	"github.com/dhth/ecsv/internal/aws"
	"github.com/dhth/ecsv/internal/types"
	"github.com/dhth/ecsv/internal/ui"
)

func process(systems []types.System, config ui.Config, awsConfigs map[string]aws.Config, maxConcFetches int) error {
	results := make(map[string]map[string]types.SystemResult)
	resultChannel := make(chan types.SystemResult)

	semaphore := make(chan struct{}, maxConcFetches)
	var wg sync.WaitGroup

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

		wg.Add(1)

		go func(system types.System) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() {
				<-semaphore
			}()
			resultChannel <- aws.FetchSystemVersion(system, awsConfig)
		}(s)
	}

	go func() {
		wg.Wait()
		close(resultChannel)
	}()

	for r := range resultChannel {
		results[r.SystemKey][r.Env] = r
	}

	output, err := ui.GetOutput(config, results)
	if err != nil {
		return err
	}
	fmt.Print(output)
	return nil
}
