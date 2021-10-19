package alarm

import (
	"media_crawling/config"

	"github.com/slack-go/slack"
)

// SlackBot : Slack Client
var SlackBot *slack.Client

func init() {
	SlackBot = slack.New(config.Secret["slack:token"])
}

// PostMessage : 입력 받은 채널로 메세지를 송부합니다.
func PostMessage(channelID string, message string) {
	if channelID == "default" {
		channelID = config.Secret["slack:channel_id"]
	}

	SlackBot.PostMessage(
		channelID,
		slack.MsgOptionText(message, true),
	)
}
