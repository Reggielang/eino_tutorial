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
	//template model
	chatTemplate := prompt.FromMessages(
		schema.FString,
		schema.SystemMessage(`你是一个{role}`),
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
	//2. 创建Chain chainTemplate => chatModel
	// 输入 map[string]any => 输出: *schema.Message
	chain := compose.NewChain[map[string]any, *schema.Message]()
	chain.
		AppendChatTemplate(chatTemplate). // 第一步：格式化模板
		AppendChatModel(chatModel)        // 第二步：调用模型

	// 3. 编译Chain
	runnable, err := chain.Compile(ctx)
	if err != nil {
		log.Fatalf("failed to compile chain: %v", err)
	}
	// 4. 执行Chain
	input := map[string]any{
		"role":  "逻辑推理专家",
		"query": `什么是Go语音？`,
	}
	output, err := runnable.Invoke(ctx, input)
	if err != nil {
		log.Fatalf("failed to run chain: %v", err)
	}
	fmt.Println("回答:", output.Content)
}
