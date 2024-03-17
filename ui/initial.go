package ui

func newModel(envSeq []string, systems []System, outFormat OutFormat) model {

	resultMap := make(map[string]map[string]string)
	var systemNames []string

	awsConfigs := make(map[string]AWSConfig)

	seenSystems := make(map[string]bool)
	seenConfigs := make(map[string]bool)

	for _, system := range systems {
		if !seenSystems[system.Key] {
			systemNames = append(systemNames, system.Key)
			seenSystems[system.Key] = true
		}
		if resultMap[system.Key] == nil {
			resultMap[system.Key] = make(map[string]string)
		}
		resultMap[system.Key][system.Env] = "..."

		if !seenConfigs[getAWSConfigKey(system)] {
			cfg, err := getAWSConfig(system)
			awsConfigs[getAWSConfigKey(system)] = AWSConfig{cfg, err}
			seenSystems[system.Key] = true
		}
	}

	errors := make([]error, 0)

	return model{
		outFormat:       outFormat,
		results:         resultMap,
		envSequence:     envSeq,
		systems:         systems,
		systemNames:     systemNames,
		numResultsToGet: len(systems),
		awsConfigs:      awsConfigs,
		printWhenReady:  true,
		errors:          errors,
	}
}
