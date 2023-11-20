package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sashabaranov/go-openai"
)

var tools = []openai.Tool{
	{
		Type: openai.ToolTypeFunction,
		Function: openai.FunctionDefinition{
			Name: "search-gifs",
			Parameters: json.RawMessage(`{
					"type": "object",
					"properties": {
						"query": {
							"type": "string",
							"description": "The query to search for e.g. \"cat\""
						}
					},
					"required": ["query"]
			}`),
		},
	},
}

var toolHandlers = map[string]func(args map[string]interface{}) string{
	"search-gifs": func(args map[string]interface{}) string {
		query := args["query"].(string)

		req, err := http.NewRequest("GET", "https://tenor.googleapis.com/v2/search", nil)
		if err != nil {
			fmt.Println("Error creating request: ", err)
			return ""
		}

		q := req.URL.Query()
		q.Add("q", query)
		q.Add("key", os.Getenv("TENOR_KEY"))
		q.Add("client_key", "my_test_app")
		q.Add("limit", "1")
		req.URL.RawQuery = q.Encode()

		resp, err := http.Get(req.URL.String())
		if err != nil {
			fmt.Println("Error getting response: ", err)
			return ""
		}

		var data map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			fmt.Println("Error decoding response: ", err)
			return ""
		}

		results := data["results"].([]interface{})
		if len(results) > 0 {
			url := results[0].(map[string]interface{})["url"].(string)
			return url
		}

		return ""
	},
}

func OpenAICall(model Model, recentMessages []*discordgo.Message) string {
	var client = openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	lastMessageContent := recentMessages[len(recentMessages)-1].Content

	moderation, err := client.Moderations(
		context.Background(),
		openai.ModerationRequest{
			Input: lastMessageContent,
		},
	)
	if err != nil {
		fmt.Println("Error calling OpenAI: ", err)
		return ""
	}

	if moderation.Results[0].Categories.Sexual {
		fmt.Println("Human message flagged: ", lastMessageContent)
		return "> :warning: **This is a warning. Do not continue asking questions like this or your access to the AI will be revoked.**"
	}

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
			Tools:    tools,
			Stop:     []string{"\n"},
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
				Tools:    tools,
			},
		)
		if err != nil {
			fmt.Println("Error calling OpenAI: ", err)
			return ""
		}
	}

	// check if any tool was used
	if resp.Choices[0].Message.ToolCalls != nil {
		for _, toolCall := range resp.Choices[0].Message.ToolCalls {
			// check if the tool exists
			if toolHandlers[toolCall.Function.Name] != nil {
				// toolCall.Function.Arguments is a string, so we need to convert it to a map
				var args map[string]interface{}
				err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
				if err != nil {
					fmt.Println("Error decoding tool arguments: ", err)
					return ""
				}

				handlerresp := toolHandlers[toolCall.Function.Name](args)
				if handlerresp != "" {
					messages = append(messages, openai.ChatCompletionMessage{Content: handlerresp, ToolCallID: toolCall.ID, Role: "tool"})
				}
			}
		}

		resp, err = client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model:    model.Identifier,
				Messages: messages,
				Tools:    tools,
			},
		)
	}

	moderation, err = client.Moderations(
		context.Background(),
		openai.ModerationRequest{
			Input: resp.Choices[0].Message.Content,
		},
	)
	if err != nil {
		fmt.Println("Error calling OpenAI: ", err)
		return ""
	}

	if moderation.Results[0].Flagged {
		fmt.Println("AI Response flagged: ", resp.Choices[0].Message.Content)
		return "> :warning: **This is a warning. Do not continue asking questions like this or your access to the AI will be revoked.** ||" + resp.Choices[0].Message.Content + "||"
	}

	resp.Choices[0].Message.Content = strings.TrimPrefix(resp.Choices[0].Message.Content, model.Prefix)

	// if there is a new line, only return the first line
	if strings.Contains(resp.Choices[0].Message.Content, "\n") {
		return strings.Split(resp.Choices[0].Message.Content, "\n")[0]
	}

	return resp.Choices[0].Message.Content
}
