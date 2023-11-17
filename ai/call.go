package ai

import (
	"github.com/bwmarrin/discordgo"
)

func Call(model Model, messages []*discordgo.Message) string {
	return model.Platform.Handler(model, messages)
}
