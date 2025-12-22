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
	//agent1:技术调研
	techAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "TechResearcher",
		Description: "负责技术调研",
		Instruction: "你是一个技术研究员，请调研相关技术方案",
		Model:       chatModel,
		OutputKey:   "tech_research", //将输出存储倒 session 的 ”tech_research“ 键
	})
	if err != nil {
		log.Fatalf("failed to create techAgent: %v", err)
	}
	//agent2:市场分析
	marketAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "MarketAnalyst",
		Description: "负责市场分析",
		Instruction: `你是一个市场分析师，请分析市场趋势和竞争对手`,
		Model:       chatModel,
		OutputKey:   "market_analysis", //将输出存储倒 session 的 ”market_analysis“ 键
	})
	if err != nil {
		log.Fatalf("failed to create market_analysis agent: %v", err)
	}
	//3. 创建 ParallelAgent
	parallelAgent, err := adk.NewParallelAgent(ctx, &adk.ParallelAgentConfig{
		Name:        "DataCollectionAgent",
		Description: "并发信息收集Agent,同时进行技术调研和市场分析",
		SubAgents:   []adk.Agent{techAgent, marketAgent},
	})
	if err != nil {
		log.Fatalf("failed to create parallelAgent: %v", err)
	}
	//4. 创建 runner
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           parallelAgent,
		EnableStreaming: false,
	})
	//5. 运行
	query := "请帮我分析以下开发一个 AI 代码助手项目的可行性"
	fmt.Printf("Query: %s\n", query)
	iter := runner.Query(ctx, query)
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
				fmt.Printf("[%s]: %s\n", event.AgentName, msg.Content)
			}
		}
	}
}
