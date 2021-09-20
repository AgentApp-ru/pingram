package config

import (
	"flag"
	"github.com/BurntSushi/toml"
)

type Config struct {
	Env                            string `toml:"env"`
	FileName                       string `toml:"filename_with_domains"`
	RedisUrl                       string `toml:"redis_url"`
	TimeoutBetweenChecks           int    `toml:"timeout_between_checks_in_seconds"`
	TimeoutBetweenHttpChecks       int    `toml:"timeout_between_http_checks_in_seconds"`
	DowntimeWithoutAlertsInSeconds int    `toml:"downtime_without_alerts_in_seconds"`
	BindAddr                       string `toml:"apiserver_port"`
	WorkingHoursStart              int    `toml:"working_hours_start"`
	WorkingHoursEnd                int    `toml:"working_hours_end"`
	SendToChat                     bool   `toml:"send_to_chat"`
	ElasticDomain                  string `toml:"elastic_domain"`
	ElasticUsername                string `toml:"elastic_username"`
	ElasticPassword                string `toml:"elastic_password"`
}

var (
	Settings   *Config
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/config.toml", "path to config file")
	flag.Parse()

	Settings = new(Config)
	_, err := toml.DecodeFile(configPath, Settings)
	if err != nil {
		panic(err)
	}
}
