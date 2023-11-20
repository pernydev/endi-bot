package ai

import (
	"endi/redacted"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func Call(model Model, messages []*discordgo.Message) string {
	// check if the model has a platform handler
	if model.Platform.Handler != nil {
		resp := model.Platform.Handler(model, messages)
		for _, word := range redacted.IllegalWords {
			resp = strings.ReplaceAll(resp, word, "****")
		}
		return resp
	}
	return "> :warning: **Invalid model**"
}
