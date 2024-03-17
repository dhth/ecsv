package ui

func newModel(envSeq []string, systems []System, outFormat OutFormat, awsConfigSource AWSConfigSource) model {

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
	}

	switch awsConfigSource {
	case SharedCfgProfileSrc:
		for _, system := range systems {
			if !seenConfigs[getSharedProfileCfgKey(&system)] {
				cfg, err := getAWSConfig(system.AWSProfile, system.AWSRegion)
				awsConfigs[getSharedProfileCfgKey(&system)] = AWSConfig{cfg, err}
				seenSystems[system.Key] = true
			}
		}
	case DefaultCfg:
		for _, system := range systems {
			switch system.IAMRoleToAssume {
			case "":
				if !seenConfigs[getDefaultCfgKey(&system)] {
					cfg, err := getDefaultConfig(system.AWSRegion)
					awsConfigs[getDefaultCfgKey(&system)] = AWSConfig{cfg, err}
				}
			default:
				if !seenConfigs[getRoleCfgKey(&system)] {
					cfg, err := getRoleConfig(system.IAMRoleToAssume, system.AWSRegion)
					awsConfigs[getRoleCfgKey(&system)] = AWSConfig{cfg, err}
				}
			}
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
		awsConfigSource: awsConfigSource,
		awsConfigs:      awsConfigs,
		printWhenReady:  true,
		errors:          errors,
	}
}
