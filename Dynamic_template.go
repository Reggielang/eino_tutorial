package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino-ext/components/model/qwen"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
	"log"
)

type ConversationStyle string

const (
	StyleProfessional ConversationStyle = "professional"
	StyleCasual       ConversationStyle = "casual"
	StyleFriendly     ConversationStyle = "friendly"
	StyleFormal       ConversationStyle = "formal"
)

func createDynamicTemplate(style ConversationStyle, domain string) prompt.ChatTemplate {
	var systemPrompt string
	switch style {
	case StyleProfessional:
		systemPrompt = fmt.Sprintf(`你是一个专业的%s专家，请根据用户的问题提供专业的回答。`, domain)
	case StyleCasual:
		systemPrompt = fmt.Sprintf(`你是一个随意的%s专家，请根据用户的问题提供随意的回答。`, domain)
	case StyleFriendly:
		systemPrompt = fmt.Sprintf(`你是一个友好的%s专家，请根据用户的问题提供友好的回答。`, domain)
	case StyleFormal:
		systemPrompt = fmt.Sprintf(`你是一个正式的%s专家，请根据用户的问题提供正式的回答。`, domain)
	default:
		systemPrompt = fmt.Sprintf(`你是一个%s专家，请根据用户的问题提供回答。`, domain)
	}
	return prompt.FromMessages(
		schema.FString,
		schema.SystemMessage(systemPrompt),
		schema.UserMessage("{query}"),
	)
}

func main() {
	ctx := context.Background()

	//1.创建调用模型
	chatModel, err := qwen.NewChatModel(ctx, &qwen.ChatModelConfig{
		APIKey:  "sk-525208983ed54ff993fcef91bd03a88a",
		Model:   "qwen-plus",
		BaseURL: "https://dashscope.aliyuncs.com/compatible-mode/v1",
	})
	if err != nil {
		log.Fatalf("failed to create chat model: %v", err)
	}
	query := `什么是微服务架构？`
	styles := []ConversationStyle{StyleProfessional, StyleCasual, StyleFriendly, StyleFormal}
	for _, style := range styles {
		fmt.Printf("==========%s 风格================", style)
		template := createDynamicTemplate(style, "微服务架构")
		messages, err := template.Format(ctx, map[string]any{
			"query": query,
		})
		if err != nil {
			log.Fatalf("failed to format messages: %v", err)
		}
		response, err := chatModel.Generate(ctx, messages)
		if err != nil {
			log.Fatalf("failed to generate response: %v", err)
		}
		fmt.Println("dynamic template result:", response.Content)
	}

}
