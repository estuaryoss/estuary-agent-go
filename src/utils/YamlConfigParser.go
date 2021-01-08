package utils

import (
	"errors"
	"github.com/estuaryoss/estuary-agent-go/src/models"
)

type YamlConfigParser struct {
	commandsList []string
}

func NewYamlConfigParser() *YamlConfigParser {
	return &YamlConfigParser{commandsList: []string{}}
}

func (configParser *YamlConfigParser) GetCommandsList(config *models.YamlConfig) []string {
	configParser.commandsList = append(configParser.commandsList, config.GetBeforeInstall()...)
	configParser.commandsList = append(configParser.commandsList, config.GetInstall()...)
	configParser.commandsList = append(configParser.commandsList, config.GetAfterInstall()...)
	configParser.commandsList = append(configParser.commandsList, config.GetBeforeScript()...)
	configParser.commandsList = append(configParser.commandsList, config.GetScript()...)
	configParser.commandsList = append(configParser.commandsList, config.GetAfterScript()...)

	return configParser.commandsList
}

func (configParser *YamlConfigParser) CheckConfig(config *models.YamlConfig) error {
	if len(config.GetScript()) == 0 {
		return errors.New("Mandatory section 'script' was not found or it was empty.")
	}
	return nil
}
