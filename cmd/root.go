package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dhth/ecsv/internal/awshelpers"
	"github.com/dhth/ecsv/internal/types"
	"github.com/dhth/ecsv/ui"
)

const (
	helpText = `Quickly check the code versions of containers running in your ECS services across various environments.

Usage: ecsv [flags]`
)

var (
	configFileName   = "ecsv/ecsv.yml"
	format           = flag.String("f", "default", "output format to use; available values: default, plaintext, html")
	htmlTemplateFile = flag.String("t", "", "path of the HTML template file to use")
)

var (
	errConfigFileFlagEmpty     = errors.New("config file flag cannot be empty")
	errCouldntGetHomeDir       = errors.New("couldn't get your home directory")
	errCouldntGetConfigDir     = errors.New("couldn't get your default config directory")
	errConfigFileExtIncorrect  = errors.New("config file must be a YAML file")
	errConfigFileDoesntExist   = errors.New("config file does not exist")
	errCouldntReadConfigFile   = errors.New("couldn't read config file")
	errCouldntParseConfigFile  = errors.New("couldn't parse config file")
	errTemplateFileDoesntExit  = errors.New("template file doesn't exist")
	errCouldntReadTemplateFile = errors.New("couldn't read template file")
	errIncorrectFormatProvided = errors.New("incorrect value for format provided")
	errEnvNotInEnvSequence     = errors.New("env not present in env-sequence")
	errNoSystemsFound          = errors.New("no systems found")
)

func Execute() error {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("%w: %s", errCouldntGetHomeDir, err.Error())
	}

	defaultConfigDir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("%w: %s", errCouldntGetConfigDir, err.Error())
	}
	defaultConfigFilePath := filepath.Join(defaultConfigDir, configFileName)

	configFilePath := flag.String("c", defaultConfigFilePath, "path of the config file")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n\nFlags:\n", helpText)
		flag.PrintDefaults()
	}

	flag.Parse()

	var outFormat types.OutputFmt
	if *format != "" {
		switch *format {
		case "default":
			outFormat = types.DefaultFmt
		case "table":
			outFormat = types.TabularFmt
		case "html":
			outFormat = types.HTMLFmt
		default:
			return fmt.Errorf("%w", errIncorrectFormatProvided)
		}
	}

	var htmlTemplate string
	if *htmlTemplateFile != "" {
		_, err := os.Stat(*htmlTemplateFile)
		if os.IsNotExist(err) {
			return fmt.Errorf("%w: path: %s", errTemplateFileDoesntExit, *htmlTemplateFile)
		}
		templateFileContents, err := os.ReadFile(*htmlTemplateFile)
		if err != nil {
			return fmt.Errorf("%w: %s", errCouldntReadTemplateFile, err.Error())
		}
		htmlTemplate = string(templateFileContents)
	}

	if *configFilePath == "" {
		return fmt.Errorf("%w", errConfigFileFlagEmpty)
	}

	configPathFull := expandTilde(*configFilePath, userHomeDir)

	if filepath.Ext(configPathFull) != ".yml" && filepath.Ext(configPathFull) != ".yaml" {
		return errConfigFileExtIncorrect
	}

	_, err = os.Stat(configPathFull)
	if os.IsNotExist(err) {
		return fmt.Errorf("%w: %s", errConfigFileDoesntExist, err.Error())
	}

	configBytes, err := os.ReadFile(configPathFull)
	if err != nil {
		return fmt.Errorf("%w: %s", errCouldntReadConfigFile, err.Error())
	}

	envSequence, systems, err := readConfig(configBytes)
	if err != nil {
		return fmt.Errorf("%w: %s", errCouldntParseConfigFile, err.Error())
	}

	// assert that all envs are present in env-sequence
	seqMap := make(map[string]bool)
	for _, s := range envSequence {
		seqMap[s] = true
	}
	for _, s := range systems {
		if !seqMap[s.Env] {
			return fmt.Errorf("%w: %s", errEnvNotInEnvSequence, s.Env)
		}
	}

	if len(systems) == 0 {
		return fmt.Errorf("%w", errNoSystemsFound)
	}

	awsConfigs := make(map[string]awshelpers.Config)

	seenSystems := make(map[string]bool)
	var systemKeys []string
	seenConfigs := make(map[string]bool)

	for _, system := range systems {
		if !seenSystems[system.Key] {
			systemKeys = append(systemKeys, system.Key)
			seenSystems[system.Key] = true
		}

		if !seenConfigs[system.AWSConfigKey()] {
			cfg, err := awshelpers.GetAWSConfig(system)
			awsConfigs[system.AWSConfigKey()] = awshelpers.Config{
				Config: cfg,
				Err:    err,
			}
			seenSystems[system.Key] = true
		}
	}

	config := ui.Config{
		EnvSequence:  envSequence,
		SystemKeys:   systemKeys,
		OutputFmt:    outFormat,
		HTMLTemplate: htmlTemplate,
	}

	err = render(systems, config, awsConfigs)
	if err != nil {
		return err
	}

	return nil
}
