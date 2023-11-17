package main

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"

	"fmt"
	"os"

	"github.com/redis/go-redis/v9"

	"endi/ai"
	"endi/global"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	global.Discord, err = discordgo.New("Bot " + os.Getenv("TOKEN"))
	if err != nil {
		fmt.Println("Error creating global.Discord session: ", err)
		return
	}

	opts, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		fmt.Println("Error parsing redis url: ", err)
		return
	}

	global.RedisC = redis.NewClient(opts)

	// models should be a []*discordgo.ApplicationCommandOptionChoice but ai.global.Models is a map[string]ai.global.Model
	models := make([]*discordgo.ApplicationCommandOptionChoice, len(ai.Models))
	i := 0
	for k, v := range ai.Models {
		models[i] = &discordgo.ApplicationCommandOptionChoice{
			Name:  v.Path + " (" + v.Author + ")",
			Value: k,
		}
		fmt.Println("global.Model loaded: ", k)
		i++
	}

	// check if the model is set
	modelpath, err := global.RedisC.Get(global.Ctx, "active-model").Result()
	if err != nil {
		fmt.Println("Error getting active model: ", err)
	} else if modelpath == "" {
		fmt.Println("No active model found")
	} else {
		global.Model = ai.Models[modelpath]
		fmt.Println("Active model: ", modelpath)
	}

	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "yhdistä-minecraft-tili",
			Description: "Yhdista Minecraft-tili global.Discord tiliisi",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "ign",
					Description: "Minecraft käyttäjänimesi",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
			},
			Type: discordgo.ChatApplicationCommand,
		},
		{
			Name: "toggle AI beta",
			Type: discordgo.UserApplicationCommand,
		},
		{
			Name:        "aseta-modeeli",
			Description: "Aseta käytettävä modeeli",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "modeeli",
					Description: "Aseta käytettävä modeeli",
					Type:        discordgo.ApplicationCommandOptionString,
					Choices:     models,
					Required:    true,
				},
			},
			Type: discordgo.ChatApplicationCommand,
		},
	}

	commandHandlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"yhdistä-minecraft-tili": linkMinecraftAccount,
		"toggle AI beta":         toggleAIBeta,
		"aseta-modeeli":          setModel,
	}

	global.Discord.AddHandlerOnce(func(s *discordgo.Session, r *discordgo.Ready) {
		fmt.Println("Bot is ready!")
		registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
		for i, v := range commands {
			fmt.Println("Registering command: ", v.Name)
			cmd, err := global.Discord.ApplicationCommandCreate(
				global.Discord.State.User.ID,
				"898265017927995422",
				v,
			)
			if err != nil {
				log.Panicf("Cannot create '%v' command: %v", v.Name, err)
			}
			registeredCommands[i] = cmd
		}
		global.Discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			fmt.Println("Command received: ", i.ApplicationCommandData().Name)
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				fmt.Println("Command hadnling: ", i.ApplicationCommandData().Name)
				h(s, i)
				fmt.Println("Command handled: ", i.ApplicationCommandData().Name)
			} else {
				fmt.Println("Command not handled: ", i.ApplicationCommandData().Name)
			}
		})
	})

	global.Discord.AddHandler(func(s *discordgo.Session, i *discordgo.MessageCreate) {
		if i.Author.Bot {
			return
		}

		if len(i.Mentions) == 0 || i.Mentions[0].ID != os.Getenv("BOT_ID") {
			return
		}

		if global.Model.Path == "" {
			global.Discord.ChannelMessageSend(i.ChannelID, "> :warning: **No active model**")
			return
		}

		// check if user has access to the AI beta
		access, err := global.RedisC.Get(global.Ctx, "ai-beta-access:"+i.Author.ID).Result()
		if err != nil || access != "true" {
			global.RedisC.Set(global.Ctx, "ai-beta-access:"+i.Author.ID, "false", 0)
			global.Discord.ChannelMessageSend(i.ChannelID, "# ||***NOT�YET***||")
			return
		}
		// get the past 5 messages
		messages, err := global.Discord.ChannelMessages(i.ChannelID, 5, i.ID, "", "")
		if err != nil {
			fmt.Println("Error getting messages: ", err)
			return
		}

		fmt.Println("Calling AI")
		// get the response from the AI
		response := ai.Call(global.Model, messages)

		// send the response
		_, err = global.Discord.ChannelMessageSend(i.ChannelID, response)
		if err != nil {
			fmt.Println("Error sending message: ", err)
		}
	})

	err = global.Discord.Open()
	if err != nil {
		fmt.Println("Error opening connection: ", err)
		return
	}

	defer global.Discord.Close()

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

	global.RedisC.Set(global.Ctx, "ign:"+i.Member.User.ID, ign, 0)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Minecraft-tili yhdistetty!",
		},
	})
}

func toggleAIBeta(s *discordgo.Session, i *discordgo.InteractionCreate) {
	currentAccess, err := global.RedisC.Get(global.Ctx, "ai-beta-access:"+i.Member.User.ID).Result()
	str := "asetettu"
	if err != nil {
		global.RedisC.Set(global.Ctx, "ai-beta-access:"+i.Member.User.ID, "true", 0)
		str = "saalittu"
	} else if currentAccess == "true" {
		global.RedisC.Set(global.Ctx, "ai-beta-access:"+i.Member.User.ID, "false", 0)
		str = "estetty"
	} else {
		global.RedisC.Set(global.Ctx, "ai-beta-access:"+i.Member.User.ID, "true", 0)
		str = "saalittu"
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "AI betan käyttö " + str + "!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

func setModel(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options

	// Or convert the slice into a map
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	modelpath := optionMap["modeeli"].StringValue()

	fmt.Println(global.Model)

	//lowercasify ign
	global.Model = ai.Models[modelpath]

	global.RedisC.Set(global.Ctx, "active-model", modelpath, 0)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "global.Model asetettu!",
		},
	})
}
