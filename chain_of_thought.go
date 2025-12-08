package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino-ext/components/model/qwen"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
	"log"
)

func main() {
	ctx := context.Background()
	// chain of thought sample
	template := prompt.FromMessages(
		schema.FString,
		schema.SystemMessage(`你是一个逻辑推理专家。	
		请按以下步骤解决问题:
		1.理解问题:复述问题的要求
		2.分析:列出解决问题需要的步骤
		3.计算:逐步执行计算
		4.验证:检查答案是否合理
		5.结论:给出最终答案`),
		schema.UserMessage("{query}"),
	)
	//1.创建调用模型
	chatModel, err := qwen.NewChatModel(ctx, &qwen.ChatModelConfig{
		APIKey:  "sk-525208983ed54ff993fcef91bd03a88a",
		Model:   "qwen-plus",
		BaseURL: "https://dashscope.aliyuncs.com/compatible-mode/v1",
	})
	if err != nil {
		log.Fatalf("failed to create chat model: %v", err)
	}
	query := `一个GO程序启动了100个goroutine,每个goroutine需要处理1000个任务。如果每个任务平均耗时10ms,且goroutine之间
					没有依赖关系，假设系统有4个CPU核心，请估算完成所有任务需要多长时间？`
	messages, err := template.Format(ctx, map[string]any{
		"query": query,
	})
	if err != nil {
		log.Fatalf("failed to format messages: %v", err)
	}
	response, err := chatModel.Generate(ctx, messages)
	if err != nil {
		log.Fatalf("failed to generate response: %v", err)
	}

	fmt.Println("chain of thought result:", response.Content)

}
