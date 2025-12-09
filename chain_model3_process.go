package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino-ext/components/model/qwen"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"log"
	"strings"
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

	//构建处理链
	chain := compose.NewChain[string, string]()

	chain.
		//step1 数据清洗
		AppendLambda(compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
			//去掉多余空格，换行等
			cleaned := strings.TrimSpace(input)
			cleaned = strings.ReplaceAll(input, "\n", "")
			return cleaned, nil
		})).
		//step2 转换为AI 分析输入
		AppendLambda(compose.InvokableLambda(func(ctx context.Context, input string) (map[string]any, error) {
			return map[string]any{"query": input}, nil
		})).
		//step3 调用模型
		AppendGraph(func() *compose.Chain[map[string]any, *schema.Message] {
			analysisChain := compose.NewChain[map[string]any, *schema.Message]()

			template := prompt.FromMessages(
				schema.FString,
				schema.SystemMessage(`你是一个文本分析专家。请分析以下文本的关键信息、主题和情感。`),
				schema.UserMessage("{query}"),
			)
			analysisChain.
				AppendChatTemplate(template). // 第一步：格式化模板
				AppendChatModel(chatModel)    // 第二步：调用模型

			return analysisChain
		}()).

		//step4 处理模型输出
		AppendLambda(compose.InvokableLambda(func(ctx context.Context, input *schema.Message) (string, error) {
			return input.Content, nil
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
	fmt.Println("chain result:", output)
}
