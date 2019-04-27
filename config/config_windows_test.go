package config

const ruleFilesConfigFile = "testdata/rules_abs_path_windows.good.yml"

var ruleFilesExpectedConf = &Config{GlobalConfig: DefaultGlobalConfig, RuleFiles: []string{"testdata\\first.rules", "testdata\\rules\\second.rules", "c:\\absolute\\third.rules"}, original: ""}
