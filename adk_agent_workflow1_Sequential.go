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
	//agent1:分析用户需求
	analyzerAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "Analyzer",
		Description: "分析用户需求，提取关键信息",
		Instruction: "你是需求分析师，请分析用户的需求，提取关键信息",
		Model:       chatModel,
		OutputKey:   "analysis", //将输出存储倒 session 的 ”analysis“ 键
	})
	if err != nil {
		log.Fatalf("failed to create analyzer agent: %v", err)
	}
	//agent2:生成方案
	solutionAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "solutionGenerator",
		Description: "根据分析结果生成解决方案",
		Instruction: `你是一个解决方案生成器。请根据需求分析结果生成详细的解决方案。可以使用 {analysis} 获取需求分析结果`,
		Model:       chatModel,
	})
	if err != nil {
		log.Fatalf("failed to create solution generator agent: %v", err)
	}
	//3. 创建 SequentialAgent
	sequentialAgent, err := adk.NewSequentialAgent(ctx, &adk.SequentialAgentConfig{
		Name:        "AnalysisWorkflow",
		Description: "分析用户需求并生成解决方案的工作流",
		SubAgents:   []adk.Agent{analyzerAgent, solutionAgent},
	})
	if err != nil {
		log.Fatalf("failed to create sequential agent: %v", err)
	}
	//4. 创建 runner
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           sequentialAgent,
		EnableStreaming: false,
	})
	//5. 运行
	query := "我想开发一个智能客服系统，需要支持多轮对话和知识库检索"
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
