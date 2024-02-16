package main

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"

	"fmt"
	"os"

	"github.com/redis/go-redis/v9"

	"endi/global"
	"endi/voice"
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
	global.SupabaseC, _ = supabase.NewClient(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_KEY"), nil)

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
	}

	commandHandlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"yhdistä-minecraft-tili": linkMinecraftAccount,
		"join":                   voice.JoinVC,
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
