package ai

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sashabaranov/go-openai"
)

func OpenAICall(model Model, recentMessages []*discordgo.Message) string {
	var client = openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	var messages []openai.ChatCompletionMessage = []openai.ChatCompletionMessage{
		{Content: "per.ny: Hei tyypit, koodasin just hauskan chatti AIn tälle servulle jos joku haluu testaa :) Sen nimi on Endi.", Role: "system"},
	}
	for _, message := range recentMessages {
		if message.Author.ID == os.Getenv("BOT_ID") {
			messages = append(messages, openai.ChatCompletionMessage{Content: model.Prefix + message.Content, Role: "assistant"})
		}
		messages = append(messages, openai.ChatCompletionMessage{Content: strings.Split(message.Author.String(), "#")[0] + ": " + message.Content, Role: "user"})
	}
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    model.Identifier,
			Messages: messages,
		},
	)
	if err != nil {
		fmt.Println("Error calling OpenAI: ", err)
		return ""
	}

	// check if it doesn't begin with the prefix
	if !strings.HasPrefix(resp.Choices[0].Message.Content, model.Prefix) {
		messages = append(messages, openai.ChatCompletionMessage{Content: "API ERROR: 403 Forbidden. You do not own that username. Please use \"endi\" instead", Role: "system"})
		resp, err = client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model:    model.Identifier,
				Messages: messages,
			},
		)
		if err != nil {
			fmt.Println("Error calling OpenAI: ", err)
			return ""
		}
	}

	return strings.TrimPrefix(resp.Choices[0].Message.Content, model.Prefix)
}
