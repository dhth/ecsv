package types

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"slices"
	"sort"
	"strings"
	"time"
)

var (
	errInvalidConfigSourceProvided   = errors.New("invalid aws-system-source provided")
	errChangesOwnerIsEmpty           = errors.New("owner (under changes) is empty")
	errChangesRepoIsEmpty            = errors.New("repo (under changes)  is empty")
	errChangesBaseNotInEnvs          = errors.New("base (under changes) is not in the provided envs")
	errChangesHeadNotInEnvs          = errors.New("head (under changes) is not in the provided envs")
	errChangesIgnorePatternIncorrect = errors.New("ignore pattern (under changes) is not valid regex")
	errSystemConfigIsIncorrect       = errors.New("system config is incorrect")
)

type OutputFmt uint

const (
	DefaultFmt OutputFmt = iota
	TabularFmt
	HTMLFmt
)

func OutputFormats() []string {
	return []string{"default", "table", "html"}
}

func (f OutputFmt) String() string {
	var value string
	switch f {
	case DefaultFmt:
		value = "default"
	case HTMLFmt:
		value = "html"
	case TabularFmt:
		value = "table"
	}

	return value
}

type AWSConfigSourceType uint

const (
	DefaultCfgType AWSConfigSourceType = iota
	SharedCfgProfileType
	AssumeRoleCfgType
)

type changesConfig struct {
	Owner         string  `yaml:"owner"`
	Repo          string  `yaml:"repo"`
	Base          string  `yaml:"base"`
	Head          string  `yaml:"head"`
	IgnorePattern *string `yaml:"ignore-pattern"`
}

type ECSVConfig struct {
	EnvSequence []string `yaml:"env-sequence"`
	Systems     []struct {
		Key  string `yaml:"key"`
		Envs []struct {
			Name            string `yaml:"name"`
			AwsConfigSource string `yaml:"aws-config-source"`
			AwsRegion       string `yaml:"aws-region"`
			Cluster         string `yaml:"cluster"`
			Service         string `yaml:"service"`
			ContainerName   string `yaml:"container-name"`
		} `yaml:"envs"`
		ChangesConfig *changesConfig `yaml:"changes"`
	} `yaml:"systems"`
}

type VersionConfig struct {
	Key                 string
	Env                 string
	AWSConfigSourceType AWSConfigSourceType
	AWSConfigSource     string
	AWSRegion           string
	IAMRoleToAssume     string
	ClusterName         string
	ServiceName         string
	ContainerName       string
}

type ChangesConfig struct {
	SystemKey     string
	Owner         string
	Repo          string
	Base          string
	Head          string
	IgnorePattern *regexp.Regexp
}

type SystemsConfig struct {
	Versions []VersionConfig
	Changes  []ChangesConfig
}

func (c ECSVConfig) Parse(keyRegex *regexp.Regexp) (SystemsConfig, []error) {
	var zero SystemsConfig

	var versionConfigs []VersionConfig
	var changesConfigs []ChangesConfig
	var errors []error

	for i, system := range c.Systems {
		var systemErrors []error
		if keyRegex != nil && !keyRegex.Match([]byte(system.Key)) {
			continue
		}

		systemEnvs := make([]string, len(system.Envs))
		for j, env := range system.Envs {
			systemEnvs[j] = env.Name
			var awsConfigType AWSConfigSourceType
			var awsConfigSource string
			switch {
			case env.AwsConfigSource == "default":
				awsConfigType = DefaultCfgType
			case strings.HasPrefix(env.AwsConfigSource, "profile:::"):
				configElements := strings.Split(env.AwsConfigSource, "profile:::")
				awsConfigSource = os.ExpandEnv(configElements[len(configElements)-1])
				awsConfigType = SharedCfgProfileType
			case strings.HasPrefix(env.AwsConfigSource, "assume-role:::"):
				configElements := strings.Split(env.AwsConfigSource, "assume-role:::")
				awsConfigSource = os.ExpandEnv(configElements[len(configElements)-1])
				awsConfigType = AssumeRoleCfgType
			default:
				systemErrors = append(systemErrors, errInvalidConfigSourceProvided)
			}

			if len(systemErrors) == 0 {
				versionConfigs = append(versionConfigs, VersionConfig{
					Key:                 system.Key,
					Env:                 env.Name,
					AWSConfigSourceType: awsConfigType,
					AWSConfigSource:     awsConfigSource,
					AWSRegion:           env.AwsRegion,
					ClusterName:         env.Cluster,
					ServiceName:         env.Service,
					ContainerName:       env.ContainerName,
				})
			}
		}

		if system.ChangesConfig != nil {
			if strings.TrimSpace(system.ChangesConfig.Owner) == "" {
				systemErrors = append(systemErrors, errChangesOwnerIsEmpty)
			}

			if strings.TrimSpace(system.ChangesConfig.Repo) == "" {
				systemErrors = append(systemErrors, errChangesRepoIsEmpty)
			}

			if !slices.Contains(systemEnvs, system.ChangesConfig.Base) {
				systemErrors = append(systemErrors, errChangesBaseNotInEnvs)
			}

			if !slices.Contains(systemEnvs, system.ChangesConfig.Head) {
				systemErrors = append(systemErrors, errChangesHeadNotInEnvs)
			}

			var ignorePattern *regexp.Regexp
			if system.ChangesConfig.IgnorePattern != nil {
				ip, err := regexp.Compile(*system.ChangesConfig.IgnorePattern)
				if err != nil {
					systemErrors = append(systemErrors, fmt.Errorf("%w: %s", errChangesIgnorePatternIncorrect, err.Error()))
				} else {
					ignorePattern = ip
				}
			}

			if len(systemErrors) == 0 {
				changesConfigs = append(changesConfigs, ChangesConfig{
					SystemKey:     system.Key,
					Owner:         system.ChangesConfig.Owner,
					Repo:          system.ChangesConfig.Repo,
					Base:          system.ChangesConfig.Base,
					Head:          system.ChangesConfig.Head,
					IgnorePattern: ignorePattern,
				})
			}
		}

		if len(systemErrors) > 0 {
			errors = append(errors, fmt.Errorf("%w; index: %d, errors: %v", errSystemConfigIsIncorrect, i+1, systemErrors))
		}
	}

	if len(errors) > 0 {
		return zero, errors
	}

	return SystemsConfig{
		Versions: versionConfigs,
		Changes:  changesConfigs,
	}, nil
}

func (vc VersionConfig) AWSConfigKey() string {
	switch vc.AWSConfigSourceType {
	case SharedCfgProfileType, AssumeRoleCfgType:
		return vc.AWSConfigSource + ":" + vc.AWSRegion
	default:
		return vc.AWSRegion
	}
}

type VersionResult struct {
	SystemKey    string
	Env          string
	Version      string
	Found        bool
	RegisteredAt *time.Time
	Err          error
}

type ChangesResult struct {
	SystemKey string
	Commits   []Commit
	Error     error
}

type Commit struct {
	SHA        string
	Message    string
	HTMLURL    string
	Author     string
	AuthoredAt string
}

type TableStyle string

func (ts TableStyle) String() string {
	return string(ts)
}

const (
	ASCIIStyle TableStyle = "ascii"
	BlankStyle TableStyle = "blank"
	DotsStyle  TableStyle = "dots"
	SharpStyle TableStyle = "sharp"
)

var styles = map[string]TableStyle{
	string(ASCIIStyle): ASCIIStyle,
	string(BlankStyle): BlankStyle,
	string(DotsStyle):  DotsStyle,
	string(SharpStyle): SharpStyle,
}

func TableStyleStrings() []string {
	values := make([]string, 0, len(styles))
	for k := range styles {
		values = append(values, k)
	}
	sort.Strings(values)
	return values
}

func GetStyle(styleStr string) (TableStyle, bool) {
	style, ok := styles[styleStr]
	return style, ok
}
