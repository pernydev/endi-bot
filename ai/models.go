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
	Prefix       string
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
	"endi.alpha.1": {
		Path:         "endi.alpha.1",
		Author:       "per.ny",
		Contributors: []string{"All of The End Discord"},
		License:      "none",
		Platform:     Platforms["openai-chatcompletions"],
		Identifier:   "ft:gpt-3.5-turbo-0613:personal::8L9USDUh",
		Prefix:       "Endi AI: ",
	},
	"endi.alpha.2": {
		Path:         "endi.alpha.2",
		Author:       "per.ny",
		Contributors: []string{"All of The End Discord"},
		License:      "none",
		Platform:     Platforms["openai-chatcompletions"],
		Identifier:   "ft:gpt-3.5-turbo-1106:personal::8LUJTol5",
		Prefix:       "Endi AI: ",
	},
	"endi.alpha.3": {
		Path:         "endi.alpha.3",
		Author:       "per.ny",
		Contributors: []string{"All of The End Discord"},
		License:      "none",
		Platform:     Platforms["openai-chatcompletions"],
		Identifier:   "ft:gpt-3.5-turbo-1106:personal::8LWi1Gvm",
		Prefix:       "Endi AI: ",
	},
	"endi.alpha.4": {
		Path:         "endi.alpha.4",
		Author:       "per.ny",
		Contributors: []string{"All of The End Discord"},
		License:      "none",
		Platform:     Platforms["openai-chatcompletions"],
		Identifier:   "ft:gpt-3.5-turbo-1106:personal::8La6UqEs",
		Prefix:       "Endi AI: ",
	},
	"endi.beta.1": {
		Path:         "endi.beta.1",
		Author:       "per.ny",
		Contributors: []string{"All of The End Discord"},
		License:      "none",
		Platform:     Platforms["openai-chatcompletions"],
		Identifier:   "ft:gpt-3.5-turbo-0613:personal::8LngQ8bE",
		Prefix:       "endi: ",
	},
	"openai.gpt-3.5-turbo": {
		Path:         "openai.gpt-3.5-turbo",
		Author:       "OpenAI",
		Contributors: []string{"OpenAI"},
		License:      "none",
		Platform:     Platforms["openai-chatcompletions"],
		Identifier:   "gpt-3.5-turbo",
		Prefix:       "Endi: ",
	},
}
