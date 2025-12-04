package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino-ext/components/model/qwen"
	"github.com/cloudwego/eino/schema"
	"log"
)

// QWEN_API_KEY=sk-525208983ed54ff993fcef91bd03a88a
func main() {
	//1.创建上下文
	ctx := context.Background()
	//2.创建chat model
	chatMode, err := qwen.NewChatModel(ctx, &qwen.ChatModelConfig{
		APIKey:  "sk-525208983ed54ff993fcef91bd03a88a",
		Model:   "qwen-plus",
		BaseURL: "https://dashscope.aliyuncs.com/compatible-mode/v1",
	})
	if err != nil {
		log.Fatalf("failed to create chat model: %v", err)
	}
	//3.准备message
	message := []*schema.Message{
		schema.SystemMessage("你是一个AI助手"),
		schema.UserMessage("请你介绍一下 Eino 框架"),
	}
	//4.调用chat model
	response, err := chatMode.Generate(ctx, message)
	if err != nil {
		log.Fatalf("failed to chat: %v", err)
	}
	//5.打印结果
	fmt.Printf("response: %v", response.Content)

	//6. 输出token使用情况
	if response.ResponseMeta != nil && response.ResponseMeta.Usage != nil {
		fmt.Println("\n Token：使用总计")
		fmt.Printf(" 输入Token：%d\n", response.ResponseMeta.Usage.PromptTokens)
		fmt.Printf(" 输出Token：%d\n", response.ResponseMeta.Usage.CompletionTokens)
		fmt.Printf(" 总计Token：%d\n", response.ResponseMeta.Usage.TotalTokens)
	}
}
