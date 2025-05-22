package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/dhth/ecsv/internal/aws"
	"github.com/dhth/ecsv/internal/changes"
	"github.com/dhth/ecsv/internal/types"
	"github.com/dhth/ecsv/internal/ui"
	"github.com/dhth/ecsv/internal/utils"
	"github.com/google/go-github/v72/github"
	"github.com/spf13/cobra"
)

const configFileName = "ecsv/ecsv.yml"

var (
	errConfigFileNotYAML         = errors.New("config file needs to be a YAML file")
	errCouldntGetUserHomeDir     = errors.New("couldn't get your home directory")
	errCouldntGetUserConfigDir   = errors.New("couldn't get your config directory")
	errConfigFileExtIncorrect    = errors.New("config file must be a YAML file")
	errConfigFileDoesntExist     = errors.New("config file does not exist")
	errCouldntReadConfigFile     = errors.New("couldn't read config file")
	errCouldntParseConfigFile    = errors.New("couldn't parse config file")
	errTemplateFileDoesntExit    = errors.New("template file doesn't exist")
	errCouldntReadTemplateFile   = errors.New("couldn't read template file")
	errIncorrectFormatProvided   = errors.New("incorrect value for format provided")
	errNoSystemsFound            = errors.New("no systems found")
	errIncorrectStyleProvided    = errors.New("incorrect style provided")
	errIncorrectKeyRegexProvided = errors.New("incorrect key regex provided")
	errGithubAuthNotConfigured   = errors.New("couldn't set up a GitHub client")
)

func Execute() error {
	rootCmd, err := NewRootCommand()
	if err != nil {
		return err
	}

	return rootCmd.Execute()
}

func NewRootCommand() (*cobra.Command, error) {
	var (
		configPath       string
		configPathFull   string
		configBytes      []byte
		homeDir          string
		keyFilter        string
		format           string
		htmlTemplateFile string
		htmlTitle        string
		htmlTitleURL     string
		htmlOpen         bool
		tableStyleStr    string
		showRegisteredAt bool
		debug            bool
	)

	rootCmd := &cobra.Command{
		Use:          "ecsv",
		Short:        "ecsv lets you quickly check the code versions of services running on ECS across various environments",
		SilenceUsage: true,
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			if !strings.HasSuffix(configPath, ".yml") && !strings.HasSuffix(configPath, ".yaml") {
				return errConfigFileNotYAML
			}

			var err error
			configPathFull = utils.ExpandTilde(configPath, homeDir)
			configBytes, err = os.ReadFile(configPathFull)
			if err != nil {
				return fmt.Errorf("%w: %w", errCouldntReadConfigFile, err)
			}

			return nil
		},
	}

	checkCmd := &cobra.Command{
		Use:          "check",
		Short:        "gather code versions and show report",
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, _ []string) error {
			var outFormat types.OutputFmt
			if format != "" {
				switch format {
				case "default":
					outFormat = types.DefaultFmt
				case "table":
					outFormat = types.TabularFmt
				case "html":
					outFormat = types.HTMLFmt
				default:
					return fmt.Errorf("%w; possible values: %v", errIncorrectFormatProvided, types.OutputFormats())
				}
			}

			var htmlTemplate string
			if htmlTemplateFile != "" {
				_, err := os.Stat(htmlTemplateFile)
				if os.IsNotExist(err) {
					return fmt.Errorf("%w: path: %s", errTemplateFileDoesntExit, htmlTemplateFile)
				}
				templateFileContents, err := os.ReadFile(htmlTemplateFile)
				if err != nil {
					return fmt.Errorf("%w: %s", errCouldntReadTemplateFile, err.Error())
				}
				htmlTemplate = string(templateFileContents)
			}

			var keyFilterRegex *regexp.Regexp
			var err error
			if keyFilter != "" {
				keyFilterRegex, err = regexp.Compile(keyFilter)
				if err != nil {
					return fmt.Errorf("%w: %s", errIncorrectKeyRegexProvided, err.Error())
				}
			}

			if filepath.Ext(configPathFull) != ".yml" && filepath.Ext(configPathFull) != ".yaml" {
				return errConfigFileExtIncorrect
			}

			_, err = os.Stat(configPathFull)
			if os.IsNotExist(err) {
				return fmt.Errorf("%w: %s", errConfigFileDoesntExist, err.Error())
			}

			envSequence, config, err := readConfig(configBytes, keyFilterRegex)
			if err != nil {
				return fmt.Errorf("%w: %s", errCouldntParseConfigFile, err.Error())
			}

			if len(config.Versions) == 0 {
				return fmt.Errorf("%w", errNoSystemsFound)
			}

			maxConcFetches, err := getMaxConcFetches()
			if err != nil {
				return err
			}

			var ghClient *github.Client
			if len(config.Changes) > 0 {
				ghClient, err = changes.GetGHClient()
				if err != nil {
					return fmt.Errorf("%w: %w", errGithubAuthNotConfigured, err)
				}
			}

			awsConfigs := make(map[string]aws.Config)

			seenSystems := make(map[string]bool)
			var systemKeys []string
			seenConfigs := make(map[string]bool)

			for _, system := range config.Versions {
				if !seenSystems[system.Key] {
					systemKeys = append(systemKeys, system.Key)
					seenSystems[system.Key] = true
				}

				if !seenConfigs[system.AWSConfigKey()] {
					cfg, err := aws.GetConfig(system)
					awsConfigs[system.AWSConfigKey()] = aws.Config{
						Config: cfg,
						Err:    err,
					}
					seenSystems[system.Key] = true
				}
			}

			uiConfig := ui.Config{
				EnvSequence:      envSequence,
				SystemKeys:       systemKeys,
				OutputFmt:        outFormat,
				ShowRegisteredAt: showRegisteredAt,
			}
			switch outFormat {
			case types.HTMLFmt:
				uiConfig.HTMLConfig = ui.HTMLOutputConfig{
					Template: htmlTemplate,
					Title:    htmlTitle,
					TitleURL: htmlTitleURL,
					Open:     htmlOpen,
				}
			case types.TabularFmt:
				tableStyle, ok := types.GetStyle(tableStyleStr)
				if !ok {
					return fmt.Errorf("%w: potential values: %q", errIncorrectStyleProvided, types.TableStyleStrings())
				}

				uiConfig.TableConfig = ui.TableOutputConfig{
					Style: tableStyle,
				}
			}

			if debug {
				fmt.Printf(`config:
%s
`, uiConfig.String())
				return nil
			}

			return process(config, uiConfig, awsConfigs, ghClient, maxConcFetches)
		},
	}

	var err error
	homeDir, err = os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errCouldntGetUserHomeDir, err.Error())
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errCouldntGetUserConfigDir, err.Error())
	}

	defaultConfigPath := filepath.Join(configDir, configFileName)

	rootCmd.PersistentFlags().StringVarP(&configPath, "config-path", "c", defaultConfigPath, "location of ecsv's config file")

	checkCmd.Flags().StringVarP(&keyFilter, "key-filter", "k", "", "regex for filtering systems (by key)")
	checkCmd.Flags().StringVarP(&format, "format", "f", "default", fmt.Sprintf("output format to use [possible values: %s]", strings.Join(types.OutputFormats(), ", ")))
	checkCmd.Flags().StringVar(&htmlTemplateFile, "html-template-file", "", "path of the HTML template file to use")
	checkCmd.Flags().StringVar(&htmlTitle, "html-title", "ecsv", "title to be used in the html output")
	checkCmd.Flags().StringVar(&htmlTitleURL, "html-title-url", "https://github.com/dhth/ecsv", "url the title in the html output should point to")
	checkCmd.Flags().BoolVar(&htmlOpen, "html-open", true, "whether to write the html output to a temporary file and open it")
	checkCmd.Flags().StringVar(&tableStyleStr, "table-style", types.ASCIIStyle.String(), fmt.Sprintf("style to use for tabular output [possible values: %s]", strings.Join(types.TableStyleStrings(), ", ")))
	checkCmd.Flags().BoolVar(&showRegisteredAt, "show-registered-at", true, "whether to show the time when the task definition corresponding to a container was registered")
	checkCmd.Flags().BoolVar(&debug, "debug", false, "whether to show debug information without running the checks")

	rootCmd.AddCommand(checkCmd)
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	return rootCmd, nil
}
