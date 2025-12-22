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
	//agent1: 主任务解决 Agent
	mainAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "MainAgent",
		Description: "负责生成初步解决方案",
		Instruction: "你是一个问题解决专家，请根据用户问题生成详细的解决方案。如果解决方案需要改进，请说明需要改进的地方",
		Model:       chatModel,
		OutputKey:   "solution",
	})
	if err != nil {
		log.Fatalf("failed to create MainAgent: %v", err)
	}
	//agent2: 批判反馈 Agent
	critiqueAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "CritiqueAgent",
		Description: "对解决方案进行批判和反馈",
		Instruction: `你是一个质量审查专家。请审查解决方案的质量，提供改进建议。如果解决方案已经完善，请明确说明"解决方案已经完善，无需进一步改进。"
						可以使用 {solution} 获取当前的解决方案`,
		Model:     chatModel,
		OutputKey: "critique",
	})
	if err != nil {
		log.Fatalf("failed to create CritiqueAgent: %v", err)
	}
	//3. 创建 LoopAgent
	loopAgent, err := adk.NewLoopAgent(ctx, &adk.LoopAgentConfig{
		Name:          "reflectionAgent",
		Description:   "迭代烦死型智能体，通过多轮迭代优化解决方案",
		SubAgents:     []adk.Agent{mainAgent, critiqueAgent},
		MaxIterations: 5, // 设置最大迭代次数
	})
	if err != nil {
		log.Fatalf("failed to create LoopAgent: %v", err)
	}
	//4. 创建 runner
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           loopAgent,
		EnableStreaming: false,
	})
	//5. 运行
	query := "如何设计一个高性能的分布式缓存系统？"
	fmt.Printf("Query: %s\n", query)
	iter := runner.Query(ctx, query)
	iteration := 0
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
				if event.AgentName == "MainAgent" {
					iteration++
					fmt.Printf("[%s] 第%d轮: %s\n", event.AgentName, iteration, msg.Content)
				} else if event.AgentName == "CritiqueAgent" {
					fmt.Printf("[%s] 第%d轮: %s\n", event.AgentName, iteration, msg.Content)
				}
				fmt.Printf("[%s] : %s ", event.AgentName, msg.Content)
			}
		}
	}
}
