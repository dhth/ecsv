package cmd

import (
	"fmt"
	"sync"

	"github.com/dhth/ecsv/internal/aws"
	"github.com/dhth/ecsv/internal/types"
	"github.com/dhth/ecsv/internal/ui"
)

func process(systemsConfig types.SystemsConfig, uiConfig ui.VersionsUIConfig, awsConfigs map[string]aws.Config, maxConcFetches int) error {
	results := make(map[string]map[string]types.VersionResult)
	resultChannel := make(chan types.VersionResult)

	semaphore := make(chan struct{}, maxConcFetches)
	var wg sync.WaitGroup

	for _, s := range systemsConfig.Versions {
		awsConfig := awsConfigs[s.AWSConfigKey()]
		if results[s.Key] == nil {
			results[s.Key] = make(map[string]types.VersionResult)
		}
		results[s.Key][s.Env] = types.VersionResult{}

		if awsConfig.Err != nil {
			results[s.Key][s.Env] = types.VersionResult{
				SystemKey: s.Key,
				Env:       s.Env,
				Err:       awsConfig.Err,
			}
			continue
		}

		wg.Add(1)

		go func(system types.VersionConfig) {
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

	output, err := ui.GetOutput(uiConfig, results)
	if err != nil {
		return err
	}
	fmt.Print(output)
	return nil
}
