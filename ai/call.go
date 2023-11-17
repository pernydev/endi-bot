package ai

import (
	"github.com/bwmarrin/discordgo"
)

func Call(model Model, messages []*discordgo.Message) string {
	// check if the model has a platform handler
	if model.Platform.Handler != nil {
		return model.Platform.Handler(model, messages)
	}
	return "> :warning: **Invalid model**"
}
