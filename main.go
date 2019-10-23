package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.worktime/")

	viper.SetDefault("slackApiToken", "")
	viper.SetDefault("slackUserName", "")
	viper.SetDefault("slackChannel", "")

	err := viper.ReadInConfig()
	if err != nil {
		log.Println(err)
	}

	subcmd := os.Args[1]
	message := ""
	if len(os.Args) > 2 {
		message = os.Args[2]
	}

	if !(subcmd == "start" || subcmd == "end" || subcmd == "status") {
		fmt.Println("ERROR! This subcommand doesn't exists.")
		return
	}

	file, err := NewFile("events.list")
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	lastEvent, err := file.GetLastEventAction()
	if err != nil {
		log.Print(err)
	}

	if lastEvent == subcmd {
		fmt.Println("Work is already " + subcmd + "ed")
	}

	if subcmd == "start" || subcmd == "end" {
		timestamp := strconv.FormatInt(time.Now().Unix(), 10)
		file.Append(timestamp + "-" + subcmd)
		lastEvent = subcmd
	}

	events, err := file.GetEvents()
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	stats, err := NewStats(events)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	fmt.Println("Work is " + lastEvent + "ed")
	fmt.Println("Daily:  " + stats.dailyPassed + " " + stats.dailyProgress + " " + stats.dailyLeft)
	fmt.Println("Weekly: " + stats.weeklyPassed + " " + stats.weeklyProgress + " " + stats.weeklyLeft)

	if message != "" && viper.GetString("slackApiToken") != "" {
		slack, _ := NewSlack(
			viper.GetString("slackApiToken"),
			viper.GetString("slackUserName"),
			viper.GetString("slackChannel"),
		)
		slack.Notify(message)
	}
}
