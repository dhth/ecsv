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

	if *configFilePath == "" {
		die("config-file cannot be empty")
	}

	configFilePathExp := expandTilde(*configFilePath)

	_, err = os.Stat(configFilePathExp)
	if os.IsNotExist(err) {
		die(cfgErrSuggestion(fmt.Sprintf("Error: file doesn't exist at %q", configFilePathExp)))
	}

	envSequence, systems, err := readConfig(configFilePathExp)
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

	ui.RenderUI(envSequence, systems)

}