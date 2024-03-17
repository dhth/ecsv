package cmd

import (
	"fmt"
	"os"
	"os/user"

	"flag"

	"github.com/dhth/ecsv/ui"
)

func die(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

var (
	format    = flag.String("format", "", "output format to use, using this will disable TUI mode; available values: plaintext, html")
	awsCfgSrc = flag.String("aws-config-source", "shared-config-profile", "source of aws configuration; values: shared-config-profile, default")
)

func Execute() {
	currentUser, err := user.Current()
	var defaultConfigFilePath string
	if err == nil {
		defaultConfigFilePath = fmt.Sprintf("%s/.config/ecsv.yml", currentUser.HomeDir)
	}
	configFilePath := flag.String("config-file", defaultConfigFilePath, "path of the config file")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n\nFlags:\n", helpText)
		flag.PrintDefaults()
	}

	flag.Parse()

	var outFormat ui.OutFormat
	if *format != "" {
		switch *format {
		case "plaintext":
			outFormat = ui.PlainTextFmt
		case "html":
			outFormat = ui.HTMLFmt
		default:
			die("ecsv only supports the following formats: plaintext, html")
		}
	}

	if *configFilePath == "" {
		die("config-file cannot be empty")
	}

	if *awsCfgSrc == "" {
		die("aws-config-source cannot be empty")
	}

	var awsConfigSource ui.AWSConfigSource
	switch *awsCfgSrc {
	case "shared-config-profile":
		awsConfigSource = ui.SharedCfgProfileSrc
	case "default":
		awsConfigSource = ui.DefaultCfg
	default:
		die("unsupported aws-config-source value provided")
	}

	configFilePathExp := expandTilde(*configFilePath)

	_, err = os.Stat(configFilePathExp)
	if os.IsNotExist(err) {
		die(cfgErrSuggestion(fmt.Sprintf("Error: file doesn't exist at %q", configFilePathExp)))
	}

	envSequence, systems, err := readConfig(configFilePathExp, awsConfigSource)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// assert that all envs are present in env-sequence
	seqMap := make(map[string]bool)
	for _, s := range envSequence {
		seqMap[s] = true
	}
	for _, s := range systems {
		if !seqMap[s.Env] {
			die("env %q found in the config but is not present in env-sequence: %q", s.Env, envSequence)
		}
	}

	if len(systems) == 0 {
		die("No systems found in config file")
	}

	ui.RenderUI(envSequence, systems, outFormat, awsConfigSource)

}
