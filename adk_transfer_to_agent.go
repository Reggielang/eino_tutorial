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

	//2. 创建多个子 Agent
	//agent1:通用agent
	generalAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "generalAgent",
		Description: "通用助手，可以处理各种问题，也可以将任务转移给专业的agent",
		Instruction: `你是一个通用助手，你可以：
						1. 直接回答简单问题。
						2. 将复杂的技术问题转移给 TechExpert
						3. 将数学问题转移给 MathExpert`,
		Model: chatModel,
	})
	if err != nil {
		log.Fatalf("failed to create generalAgent agent: %v", err)
	}
	//agent2:技术专家
	techAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "TechExpert",
		Description: "技术专家，处理编程和技术问题",
		Instruction: `你是一个技术专家，请详细解答编程和技术相关的问题。`,
		Model:       chatModel,
	})
	if err != nil {
		log.Fatalf("failed to create TechExpert agent: %v", err)
	}
	//3. 创建 数学专家
	mathAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "mathAgent",
		Description: "数学专家，处理数学问题",
		Instruction: `你是一个数学专家，请详细解答数学相关的问题。`,
	})
	if err != nil {
		log.Fatalf("failed to create mathAgent agent: %v", err)
	}
	//4. 设置 agent 关系
	generalAgentWithSubs, err := adk.SetSubAgents(ctx, generalAgent, []adk.Agent{techAgent, mathAgent})
	if err != nil {
		log.Fatalf("failed to set sub agents: %v", err)
	}
	//5. 创建 runner
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           generalAgentWithSubs,
		EnableStreaming: false,
	})
	//5. 测试
	query := []string{
		"你好，今天天气如何？",
		"go 语言中如何实现并发？",
		"如何计算圆的面积？",
	}
	for _, q := range query {
		fmt.Printf("Query: %s\n", q)
		iter := runner.Query(ctx, q)
		for {
			event, ok := iter.Next()
			if !ok {
				break
			}
			if event.Err != nil {
				log.Fatalf("failed to run agent: %v", event.Err)
			}
			if event.Output != nil && event.Output.MessageOutput != nil {
				msg := event.Output.MessageOutput.Message
				if msg != nil {
					//检查是否有工具调用
					if len(msg.ToolCalls) > 0 {
						for _, toolCall := range msg.ToolCalls {
							if toolCall.Function.Name == "transfer_to_agent" {
								fmt.Printf("%s Transfer to agent: %s\n", event.AgentName, toolCall.Function.Arguments)
							}
						}
					} else {
						fmt.Printf("[%s]: %s\n", event.AgentName, msg.Content)
					}
				}
			}
		}
	}
}
