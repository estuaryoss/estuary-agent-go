package models

type YamlConfig struct {
	Env           map[string]string `yaml:"env"`
	BeforeInstall []string          `yaml:"before_install"`
	Install       []string          `yaml:"install"`
	AfterInstall  []string          `yaml:"after_install"`
	BeforeScript  []string          `yaml:"before_script"`
	Script        []string          `yaml:"script"`
	AfterScript   []string          `yaml:"after_script"`
}

func NewYamlConfig() *YamlConfig {
	config := &YamlConfig{}
	return config
}

func (config *YamlConfig) GetEnv() map[string]string {
	return config.Env
}

func (config *YamlConfig) SetEnv(env map[string]string) {
	config.Env = env
}

func (config *YamlConfig) GetBeforeInstall() []string {
	return config.BeforeInstall
}

func (config *YamlConfig) SetBeforeInstall(beforeInstall []string) {
	config.BeforeInstall = beforeInstall
}

func (config *YamlConfig) GetInstall() []string {
	return config.Install
}

func (config *YamlConfig) SetInstall(install []string) {
	config.Install = install
}

func (config *YamlConfig) GetAfterInstall() []string {
	return config.AfterInstall
}

func (config *YamlConfig) SetAfterInstall(afterInstall []string) {
	config.AfterInstall = afterInstall
}

func (config *YamlConfig) GetBeforeScript() []string {
	return config.BeforeScript
}

func (config *YamlConfig) SetBeforeScript(beforeScript []string) {
	config.BeforeScript = beforeScript
}

func (config *YamlConfig) GetScript() []string {
	return config.Script
}

func (config *YamlConfig) SetScript(script []string) {
	config.Script = script
}

func (config *YamlConfig) GetAfterScript() []string {
	return config.AfterScript
}

func (config *YamlConfig) SetAfterScript(afterScript []string) {
	config.AfterScript = afterScript
}
