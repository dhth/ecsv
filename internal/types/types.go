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
	errInvalidConfigSourceProvided = errors.New("invalid aws-system-source provided")
	errVcsPlatformisInvalid        = errors.New("invalid VCS platform provided")
	errChangelogRepoIsEmpty        = errors.New("repo is empty")
	errChangelogBaseNotInEnvs      = errors.New("changelog base is not in the provided envs")
	errChangelogHeadNotInEnvs      = errors.New("changelog head is not in the provided envs")
	errSystemConfigIsIncorrect     = errors.New("system config is incorrect")
)

type OutputFmt uint

const (
	DefaultFmt OutputFmt = iota
	TabularFmt
	HTMLFmt
)

type VCSPlatformType uint

const (
	Github VCSPlatformType = iota
)

func stringToPlatform(value string) (VCSPlatformType, bool) {
	var zero VCSPlatformType
	if value == "github" {
		return Github, true
	}

	return zero, false
}

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

type changeLogConfig struct {
	VCSPlatform string `yaml:"vcs-platform"`
	Repo        string `yaml:"repo"`
	Base        string `yaml:"base"`
	Head        string `yaml:"head"`
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
		ChangelogConfig *changeLogConfig `yaml:"changelog"`
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

type ChangeLogConfig struct {
	SystemKey   string
	VCSPlatform VCSPlatformType
	Repo        string
	Base        string
	Head        string
}

type SystemsConfig struct {
	Versions   []VersionConfig
	Changelogs []ChangeLogConfig
}

func (c ECSVConfig) Parse(keyRegex *regexp.Regexp) (SystemsConfig, []error) {
	var zero SystemsConfig

	var versions []VersionConfig
	var changelogs []ChangeLogConfig
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
				versions = append(versions, VersionConfig{
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

		if system.ChangelogConfig != nil {
			vcs, ok := stringToPlatform(system.ChangelogConfig.VCSPlatform)
			if !ok {
				systemErrors = append(systemErrors, errVcsPlatformisInvalid)
			}

			if strings.TrimSpace(system.ChangelogConfig.Repo) == "" {
				systemErrors = append(systemErrors, errChangelogRepoIsEmpty)
			}

			if !slices.Contains(systemEnvs, system.ChangelogConfig.Base) {
				systemErrors = append(systemErrors, errChangelogBaseNotInEnvs)
			}

			if !slices.Contains(systemEnvs, system.ChangelogConfig.Head) {
				systemErrors = append(systemErrors, errChangelogHeadNotInEnvs)
			}

			if len(systemErrors) == 0 {
				changelogs = append(changelogs, ChangeLogConfig{
					SystemKey:   system.Key,
					VCSPlatform: vcs,
					Repo:        system.ChangelogConfig.Repo,
					Base:        system.ChangelogConfig.Base,
					Head:        system.ChangelogConfig.Head,
				})
			}
		}

		if len(systemErrors) > 0 {
			errors = append(errors, fmt.Errorf("%w; index: %d, errors: %v", errSystemConfigIsIncorrect, i, systemErrors))
		}
	}

	if len(errors) > 0 {
		return zero, errors
	}

	return SystemsConfig{
		Versions:   versions,
		Changelogs: changelogs,
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
