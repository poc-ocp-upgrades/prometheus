package config

const ruleFilesConfigFile = "testdata/rules_abs_path.good.yml"

var ruleFilesExpectedConf = &Config{GlobalConfig: DefaultGlobalConfig, RuleFiles: []string{"testdata/first.rules", "testdata/rules/second.rules", "/absolute/third.rules"}, original: ""}
