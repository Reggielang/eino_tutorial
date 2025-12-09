package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino-ext/components/model/qwen"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
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

	//创建并行节点
	parallel := compose.NewParallel()

	//任务1：提取关键词
	parallel.AddLambda("keywords", compose.InvokableLambda(
		func(ctx context.Context, input map[string]any) (string, error) {
			fmt.Println("并行任务1：提取关键词")

			template := prompt.FromMessages(
				schema.FString,
				schema.SystemMessage(`你是一个关键词提取专家。请从以下文本中提取关键词。`),
				schema.UserMessage("{query}"),
			)
			messages, _ := template.Format(ctx, input)
			response, err := chatModel.Generate(ctx, messages)
			if err != nil {
				return "", err
			}
			return response.Content, nil
		},
	))
	//任务2：情感分析
	parallel.AddLambda("sentiment", compose.InvokableLambda(
		func(ctx context.Context, input map[string]any) (string, error) {
			fmt.Println("并行任务2：情感分析")

			template := prompt.FromMessages(
				schema.FString,
				schema.SystemMessage(`你是一个情感分析专家。请分析以下文本的情感。`),
				schema.UserMessage("{query}"),
			)

			messages, _ := template.Format(ctx, input)
			response, err := chatModel.Generate(ctx, messages)
			if err != nil {
				return "", err
			}
			return response.Content, nil
		}))
	//任务3：摘要生成
	parallel.AddLambda("summary", compose.InvokableLambda(
		func(ctx context.Context, input map[string]any) (string, error) {
			fmt.Println("并行任务3：摘要生成")
			template := prompt.FromMessages(
				schema.FString,
				schema.SystemMessage(`你是一个摘要生成专家。请为以下文本生成摘要。`),
				schema.UserMessage("{query}"))

			messages, _ := template.Format(ctx, input)
			response, err := chatModel.Generate(ctx, messages)
			if err != nil {
				return "", err
			}
			return response.Content, nil
		}))

	//构建处理链
	chain := compose.NewChain[string, map[string]any]()

	chain.
		//准备输入
		AppendLambda(compose.InvokableLambda(func(ctx context.Context, input string) (map[string]any, error) {
			return map[string]any{"query": input}, nil
		})).
		//并行执行三个任务
		AppendParallel(parallel).
		//合并结果
		AppendLambda(compose.InvokableLambda(func(ctx context.Context, results map[string]any) (map[string]any, error) {
			return results, nil
		}))

	runnable, err := chain.Compile(ctx)
	if err != nil {
		log.Fatalf("failed to compile chain: %v", err)
	}
	input := `
			Eino        是一个   强大的      AI  开发框架。
			它提供了丰富的组件，如模型、模板、链式处理等，可以帮助开发者快速构建强大的AI应用。
			开发者             可以 快速构建AI       应用。
             `
	output, err := runnable.Invoke(ctx, input)
	if err != nil {
		log.Fatalf("failed to invoke chain: %v", err)
	}
	fmt.Printf("关键词: %s\n", output["keywords"])
	fmt.Printf("情感: %s\n", output["sentiment"])
	fmt.Printf("摘要: %s\n", output["summary"])
}
