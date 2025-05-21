package cmd

import (
	"fmt"
	"os"
	"sort"
	"sync"

	"github.com/dhth/ecsv/internal/aws"
	"github.com/dhth/ecsv/internal/changes"
	"github.com/dhth/ecsv/internal/types"
	"github.com/dhth/ecsv/internal/ui"
	"github.com/google/go-github/v72/github"
)

func process(
	config types.Config,
	uiConfig ui.Config,
	awsConfigs map[string]aws.Config,
	maxConcFetches int,
) error {
	versionResults := make(map[string]map[string]types.VersionResult)
	resultChannel := make(chan types.VersionResult)

	semaphore := make(chan struct{}, maxConcFetches)
	var wg sync.WaitGroup

	for _, s := range config.Versions {
		awsConfig := awsConfigs[s.AWSConfigKey()]
		if versionResults[s.Key] == nil {
			versionResults[s.Key] = make(map[string]types.VersionResult)
		}
		versionResults[s.Key][s.Env] = types.VersionResult{}

		if awsConfig.Err != nil {
			versionResults[s.Key][s.Env] = types.VersionResult{
				SystemKey: s.Key,
				Env:       s.Env,
				Err:       awsConfig.Err,
			}
			continue
		}

		wg.Add(1)

		go func(system types.VersionsConfig) {
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
		versionResults[r.SystemKey][r.Env] = r
	}

	changesResultChan := make(chan types.ChangesResult)

	//nolint:prealloc
	var changesResults []types.ChangesResult

	if uiConfig.OutputFmt == types.HTMLFmt && len(config.Changes) > 0 {
		chSemaphore := make(chan struct{}, maxConcFetches)
		var changesWg sync.WaitGroup

		client := github.NewClient(nil).WithAuthToken(os.Getenv("GH_TOKEN"))
		for _, changesConfig := range config.Changes {
			vrm, ok := versionResults[changesConfig.SystemKey]

			// TODO: handle these conditions related to inconsistent state
			if !ok {
				continue
			}

			vrBase, ok := vrm[changesConfig.Base]
			if !ok {
				continue
			}

			if vrBase.Err != nil {
				continue
			}

			vrHead, ok := vrm[changesConfig.Head]
			if !ok {
				continue
			}

			if vrHead.Err != nil {
				continue
			}

			if vrBase.Version == vrHead.Version {
				continue
			}

			changesWg.Add(1)
			go func(baseRef, headRef string) {
				defer changesWg.Done()
				chSemaphore <- struct{}{}
				defer func() {
					<-chSemaphore
				}()
				changesResultChan <- changes.FetchChanges(client,
					changesConfig,
					baseRef,
					headRef)
			}(vrBase.Version,
				vrHead.Version,
			)
		}

		go func() {
			changesWg.Wait()
			close(changesResultChan)
		}()

		for r := range changesResultChan {
			changesResults = append(changesResults, r)
		}

		sort.Slice(changesResults, func(i, j int) bool {
			return changesResults[i].Config.SystemKey < changesResults[j].Config.SystemKey
		})
	}

	output, err := ui.GetOutput(uiConfig, versionResults, changesResults)
	if err != nil {
		return err
	}
	fmt.Print(output)
	return nil
}
