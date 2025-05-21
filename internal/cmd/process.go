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
	systemsConfig types.SystemsConfig,
	uiConfig ui.VersionsUIConfig,
	awsConfigs map[string]aws.Config,
	maxConcFetches int,
) error {
	versionResults := make(map[string]map[string]types.VersionResult)
	resultChannel := make(chan types.VersionResult)

	semaphore := make(chan struct{}, maxConcFetches)
	var wg sync.WaitGroup

	for _, s := range systemsConfig.Versions {
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
		versionResults[r.SystemKey][r.Env] = r
	}

	changelogResultChan := make(chan types.ChangesResult)

	//nolint:prealloc
	var changesResults []types.ChangesResult

	cLSemaphore := make(chan struct{}, maxConcFetches)
	var clWg sync.WaitGroup

	client := github.NewClient(nil).WithAuthToken(os.Getenv("GH_TOKEN"))
	for _, changesConfig := range systemsConfig.Changes {
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

		if vrBase == vrHead {
			continue
		}

		clWg.Add(1)
		go func(systemKey, owner, repo, baseRef, headRef string) {
			defer clWg.Done()
			cLSemaphore <- struct{}{}
			defer func() {
				<-cLSemaphore
			}()
			changelogResultChan <- changes.FetchChanges(client, systemKey, owner, repo, baseRef, headRef)
		}(changesConfig.SystemKey,
			changesConfig.Owner,
			changesConfig.Repo,
			vrBase.Version,
			vrHead.Version)
	}

	go func() {
		clWg.Wait()
		close(changelogResultChan)
	}()

	for r := range changelogResultChan {
		changesResults = append(changesResults, r)
	}

	sort.Slice(changesResults, func(i, j int) bool {
		return changesResults[i].SystemKey < changesResults[j].SystemKey
	})

	output, err := ui.GetOutput(uiConfig, versionResults, changesResults)
	if err != nil {
		return err
	}
	fmt.Print(output)
	return nil
}
