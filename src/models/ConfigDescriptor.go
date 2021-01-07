package models

var Description = make(map[string]interface{})

func SetDescription(config *YamlConfig, description interface{}) {
	Description["description"] = description
	Description["config"] = config
}

func GetDescription() interface{} {
	return Description
}
