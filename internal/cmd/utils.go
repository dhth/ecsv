package cmd

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

const (
	maxConcurrentFetchesDefault        = 10
	maxConcurrentFetchesUpperThreshold = 50
	maxConcurrentFetchesEnvVar         = "ECSV_MAX_CONCURRENT_FETCHES"
)

var errMaxConcFetchesIsInvalid = errors.New("maximum concurrent fetches is invalid")

func getMaxConcFetches() (int, error) {
	zero := maxConcurrentFetchesDefault
	userProvidedStr := os.Getenv(maxConcurrentFetchesEnvVar)

	if userProvidedStr == "" {
		return maxConcurrentFetchesDefault, nil
	}

	maxFetches, err := strconv.Atoi(userProvidedStr)
	if err != nil {
		return zero, fmt.Errorf("%w: %s needs to be an integer", errMaxConcFetchesIsInvalid, maxConcurrentFetchesEnvVar)
	}
	if maxFetches <= 0 || maxFetches > maxConcurrentFetchesUpperThreshold {
		return zero, fmt.Errorf("%w: value needs to be between [1, %d] (both inclusive)", errMaxConcFetchesIsInvalid, maxConcurrentFetchesUpperThreshold)
	}

	return maxFetches, nil
}
