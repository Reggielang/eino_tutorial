package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino-ext/components/model/qwen"
	"github.com/cloudwego/eino/adk"
	"log"
)

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
	//2. 创建ChatModelAgent
	chatAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "SimpleAssistant",
		Description: "A simple assistant that can answer questions",
		Instruction: "You are a helpful assistant. Please answer the user's question as accurately as possible.",
		Model:       chatModel,
		ToolsConfig: adk.ToolsConfig{}, //不使用工具
	})
	if err != nil {
		log.Fatalf("failed to create chat model agent: %v", err)
	}
	//3. 创建 Runner
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           chatAgent,
		EnableStreaming: false,
	})
	//4. 运行Agent
	query := "什么是JAVA?"
	fmt.Printf("Query: %s\n", query)

	iter := runner.Query(ctx, query)
	for {
		event, ok := iter.Next()
		if !ok {
			//迭代器关闭，退出循环
			break
		}
		if event.Err != nil {
			log.Fatalf("failed to run agent: %v", event.Err)
		}
		if event.Output != nil && event.Output.MessageOutput != nil {
			msg := event.Output.MessageOutput.Message
			if msg != nil {
				fmt.Printf("Response: %s\n", msg.Content)
			}
		}
	}
}
