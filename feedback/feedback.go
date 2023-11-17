package endi

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"endi/global"
)

func FeedbackInit() {
	global.Discord.AddHandler(func(s *discordgo.Session, i *discordgo.MessageReactionAdd) {
		fmt.Println("Reaction added")
	})
}
