package ai

import (
	"github.com/bwmarrin/discordgo"
)

type Model struct {
	Path         string
	Author       string
	Contributors []string
	License      string
	Platform     Platform
	Identifier   string
}

type Platform struct {
	Name    string
	Handler func(model Model, messages []*discordgo.Message) string
}

var Platforms = map[string]Platform{
	"openai-chatcompletions": {
		Name:    "OpenAI Chat Completions",
		Handler: OpenAICall,
	},
}

var Models = map[string]Model{
	"endi.alpha.text.4": {
		Path:         "endi",
		Author:       "per.ny",
		Contributors: []string{"All of The End Discord"},
		License:      "none",
		Platform:     Platforms["openai-chatcompletions"],
		Identifier:   "ft:gpt-3.5-turbo-1106:personal::8La6UqEs",
	},
}
