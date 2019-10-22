package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

type Slack struct {
	apiToken string
	user     string
	channel  string
}

func NewSlack(apiToken string, user string, channel string) (*Slack, error) {
	slack := &Slack{
		apiToken: apiToken,
		user:     user,
		channel:  channel,
	}

	return slack, nil
}

func (slack *Slack) Notify(message string) {
	body := fmt.Sprintf(`{
		"text":"%s",
		"channel":"%s",
		"username":"%s",
		"as_user":true
		}`, message, slack.channel, slack.user)
	auth := fmt.Sprintf("Bearer %s", slack.apiToken)
	jsonStr := []byte(body)
	req, err := http.NewRequest("POST", "https://slack.com/api/chat.postMessage", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", auth)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
}
