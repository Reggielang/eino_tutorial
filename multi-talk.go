package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/cloudwego/eino-ext/components/model/qwen"
	"github.com/cloudwego/eino/schema"
	"log"
	"os"
	"strings"
)

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
	messages := []*schema.Message{
		schema.SystemMessage("你是一个AI助手"),
	}
	//4.调用chat model
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("请输入你的问题(输出exit退出)：")
	for scanner.Scan() {
		fmt.Print("\n 你：")
		if !scanner.Scan() {
			break
		}
		userInput := strings.TrimSpace(scanner.Text())
		if userInput == "exit" {
			fmt.Println("程序退出，再见")
			break
		}
		if userInput == "" {
			continue
		}
		//添加用户信息
		messages = append(messages, schema.UserMessage(userInput))
		//调用chat model
		response, err := chatMode.Generate(ctx, messages)
		if err != nil {
			log.Fatalf("failed to chat: %v", err)
			continue
		}
		//添加AI响应到历史
		messages = append(messages, response)
		fmt.Printf("AI response: %s", response.Content)
	}
}
