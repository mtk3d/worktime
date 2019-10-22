package main

import (
	"github.com/spf13/viper"
)

type Config struct {
	slackApiToken string
	slackUserName string
	slackChannel  string
}

func NewConfig() {
	viper.AddConfigPath("$HOME/.worktime")
	viper.SetDefault("slackChannel", "CFCAZ895F")
}
