package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino-ext/components/model/qwen"
	"github.com/cloudwego/eino/adk"
	"io"
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

	//2. 创建 Agent
	agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "BookRecommender",
		Description: "书籍推荐 agent",
		Instruction: "你是一个书籍推荐专家，请根据用户需求推荐合适的书籍",
		Model:       chatModel,
	})
	if err != nil {
		log.Fatalf("failed to create techAgent: %v", err)
	}

	//3. 创建 runner
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           agent,
		EnableStreaming: true,
	})
	//5. 运行
	query := "我想学习go语言，请推荐一些书籍"
	fmt.Printf("Query: %s\n", query)
	iter := runner.Query(ctx, query)

	//6. 处理结果
	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			log.Fatalf("failed to run agent: %v", event.Err)
		}
		if event.Output != nil && event.Output.MessageOutput != nil {
			//处理流式输出
			if event.Output.MessageOutput.IsStreaming {
				stream := event.Output.MessageOutput.MessageStream
				for {
					msg, err := stream.Recv()
					if err != nil {
						if err == io.EOF {
							break
						}
						log.Fatalf("failed to receive message: %v", err)
					}
					if msg != nil && msg.Content != "" {
						fmt.Print(msg.Content)
					}
				}

			}
		}
	}
}
