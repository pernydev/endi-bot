package ai

import (
	"context"
	"fmt"
	"os"
	"regexp"
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
			messages = append(messages, openai.ChatCompletionMessage{Content: "Endi AI: " + message.Content, Role: "assistant"})
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

	// if text does not start with "Endi AI: "
	if !strings.HasPrefix(resp.Choices[0].Message.Content, "Endi AI: ") {
		fmt.Println("Error calling OpenAI: ", resp.Choices[0].Message.Content)
		return OpenAICall(model, recentMessages)
	}

	if matched, err := regexp.MatchString(`\b[a-z]+\:\s`, resp.Choices[0].Message.Content); err != nil || matched {
		fmt.Println("Error calling OpenAI: ", resp.Choices[0].Message.Content)
		return OpenAICall(model, recentMessages)
	}

	return strings.ReplaceAll(resp.Choices[0].Message.Content, "Endi AI: ", "")
}
