package main

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"

	"fmt"
	"os"

	"context"

	"github.com/redis/go-redis/v9"

	"endi/ai"
)

var (
	discord *discordgo.Session
	redisC  *redis.Client
	ctx     = context.Background()
	model   = ai.Models["endi.alpha.text.4"]
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	discord, err = discordgo.New("Bot " + os.Getenv("TOKEN"))
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	opts, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		fmt.Println("Error parsing redis url: ", err)
		return
	}

	redisC = redis.NewClient(opts)

	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "yhdistä-minecraft-tili",
			Description: "Yhdista Minecraft-tili Discord tiliisi",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "ign",
					Description: "Minecraft käyttäjänimesi",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
			},
		},
	}

	commandHandlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"yhdistä-minecraft-tili": linkMinecraftAccount,
	}

	discord.AddHandlerOnce(func(s *discordgo.Session, r *discordgo.Ready) {
		fmt.Println("Bot is ready!")
		registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
		for i, v := range commands {
			fmt.Println("Registering command: ", v.Name)
			cmd, err := discord.ApplicationCommandCreate(
				discord.State.User.ID,
				"898265017927995422",
				v,
			)
			if err != nil {
				log.Panicf("Cannot create '%v' command: %v", v.Name, err)
			}
			registeredCommands[i] = cmd
		}
		discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			fmt.Println("Command received: ", i.ApplicationCommandData().Name)
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		})
	})

	discord.AddHandler(func(s *discordgo.Session, i *discordgo.MessageCreate) {
		if len(i.Mentions) == 0 || i.Mentions[0].ID != os.Getenv("BOT_ID") || i.Author.ID == os.Getenv("BOT_ID") {
			return
		}

		if i.Author.ID != os.Getenv("OWNER_ID") {
			discord.ChannelMessageSend(i.ChannelID, "||***Not yet***||")
			return
		}

		// get the past 5 messages
		messages, err := discord.ChannelMessages(i.ChannelID, 5, i.ID, "", "")
		if err != nil {
			fmt.Println("Error getting messages: ", err)
			return
		}

		// get the response from the AI
		response := ai.Call(model, messages)

		// send the response
		_, err = discord.ChannelMessageSend(i.ChannelID, response)
		if err != nil {
			fmt.Println("Error sending message: ", err)
		}
	})

	err = discord.Open()
	if err != nil {
		fmt.Println("Error opening connection: ", err)
		return
	}

	defer discord.Close()

	fmt.Println("Bot is now running. Press CTRL-C to exit.")

	<-make(chan struct{})
}

func linkMinecraftAccount(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options

	// Or convert the slice into a map
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	ign := optionMap["ign"].StringValue()

	fmt.Println(ign)

	//lowercasify ign
	ign = strings.ToLower(ign)

	redisC.Set(ctx, "ign:"+i.Member.User.ID, ign, 0)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Minecraft-tili yhdistetty!",
		},
	})
}
