package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino-ext/components/model/qwen"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"log"
	"os"
)

type ArticleRequest struct {
	Topic    string
	Keywords []string
	Length   int //目标字数
}

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
	//2. 构建文章生成流水线
	chain := compose.NewChain[ArticleRequest, string]()

	chain.
		AppendLambda(compose.InvokableLambda(func(ctx context.Context, req ArticleRequest) (string, error) {
			fmt.Println("---生成大纲---")
			template := prompt.FromMessages(
				schema.FString,
				schema.SystemMessage("你是一个专业的内容策划师，请根据主题和关键词生成文章大纲。"),
				schema.UserMessage("主题: {topic}\n关键词: {keywords},请生成一个3-5点的文章大纲。"),
			)

			messages, _ := template.Format(ctx, map[string]any{
				"topic":    req.Topic,
				"keywords": fmt.Sprintf("%v", req.Keywords),
			})
			response, err := chatModel.Generate(ctx, messages)
			if err != nil {
				return "", err
			}
			return response.Content, nil

		})).
		AppendLambda(compose.InvokableLambda(func(ctx context.Context, outline string) (string, error) {
			fmt.Println("---生成文章---")
			template := prompt.FromMessages(
				schema.FString,
				schema.SystemMessage("你是一个专业的内容创作专家，请根据大纲生成文章。"),
				schema.UserMessage("大纲: {outline}\n请根据大纲生成一篇文章，字数控制在500字以内。"),
			)
			messages, _ := template.Format(ctx, map[string]any{
				"outline": outline,
			})
			response, err := chatModel.Generate(ctx, messages)
			if err != nil {
				return "", err
			}
			return response.Content, nil
		})).
		AppendLambda(compose.InvokableLambda(func(ctx context.Context, draft string) (string, error) {
			fmt.Println("---润色文章---")
			template := prompt.FromMessages(
				schema.FString,
				schema.SystemMessage("你是一个专业的内容编辑专家，请根据大纲润色文章。"),
				schema.UserMessage("文章: {draft}\n请根据大纲润色文章，字数控制在500字以内。"),
			)
			messages, _ := template.Format(ctx, map[string]any{
				"draft": draft,
			})
			response, err := chatModel.Generate(ctx, messages)
			if err != nil {
				return "", err
			}
			return response.Content, nil
		})).
		//格式化输出
		AppendLambda(compose.InvokableLambda(func(ctx context.Context, content string) (string, error) {
			fmt.Println("---格式化输出---")
			//添加markdown格式
			formatted := fmt.Sprintf("# 生成的文章 \n\n%s", content)
			return formatted, nil
		}))

	// 3. 编译Chain
	runnable, err := chain.Compile(ctx)
	if err != nil {
		log.Fatalf("failed to compile chain: %v", err)
	}
	// 生成文章
	request := ArticleRequest{
		Topic:    "Eino框架",
		Keywords: []string{"Eino", "AI", "开发框架", "组件", "模型", "模板", "链式处理"},
		Length:   500,
	}
	output, err := runnable.Invoke(ctx, request)
	if err != nil {
		log.Fatalf("failed to run chain: %v", err)
	}
	fmt.Println("回答:", output)

	//保存到文件
	os.WriteFile("AI_Eino.md", []byte(output), 0644)
	fmt.Println("\n 文章已保存")
}
