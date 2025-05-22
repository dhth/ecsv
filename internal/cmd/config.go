package cmd

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/dhth/ecsv/internal/types"

	"gopkg.in/yaml.v3"
)

var (
	errConfigIsInvalidYAML = errors.New("config file is not valid YAML")
	errConfigIsInvalid     = errors.New("invalid config provided")
	errEnvNotInEnvSequence = errors.New("env not present in env-sequence")
)

func readConfig(configBytes []byte, keyRegex *regexp.Regexp) ([]string, types.Config, error) {
	var zero types.Config
	ecsvConfig := types.ECSVConfig{}
	err := yaml.Unmarshal(configBytes, &ecsvConfig)
	if err != nil {
		return nil, zero, fmt.Errorf("%w: %s", errConfigIsInvalidYAML, err.Error())
	}

	config, errors := ecsvConfig.Parse(keyRegex)
	if len(errors) > 0 {
		errMsgs := make([]string, len(errors))
		for i, err := range errors {
			errMsgs[i] = fmt.Sprintf("- %s", err.Error())
		}
		return nil, zero, fmt.Errorf("%w; errors:\n%s", errConfigIsInvalid, strings.Join(errMsgs, "\n"))
	}

	// assert that all envs are present in env-sequence
	seqMap := make(map[string]bool)
	for _, s := range ecsvConfig.EnvSequence {
		seqMap[s] = true
	}

	for _, vc := range config.Versions {
		if !seqMap[vc.Env] {
			return nil, zero, fmt.Errorf("%w: %s", errEnvNotInEnvSequence, vc.Env)
		}
	}

	return ecsvConfig.EnvSequence, config, nil
}
