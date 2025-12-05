package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino-ext/components/model/qwen"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
	"log"
)

type PromptTemplate struct {
}

// 翻译助手模板
func (p *PromptTemplate) Translator(sourceLang, targetLang string) prompt.ChatTemplate {
	return prompt.FromMessages(
		schema.FString,
		schema.SystemMessage(fmt.Sprintf(
			"你是一个翻译助手，请将%s翻译为%s", sourceLang, targetLang)),
		schema.UserMessage("{text}"),
	)
}

// 代码审查模板
func (p *PromptTemplate) CodeReview(code string) prompt.ChatTemplate {
	return prompt.FromMessages(
		schema.FString,
		schema.SystemMessage("你是一个代码审查助手，请审查以下代码"),
		schema.UserMessage("请审查以下代码：```{language} \n {code} \n```"),
	)
}

func main() {
	ctx := context.Background()
	templates := &PromptTemplate{}
	//1.创建调用模型
	chatModel, err := qwen.NewChatModel(ctx, &qwen.ChatModelConfig{
		APIKey:  "sk-525208983ed54ff993fcef91bd03a88a",
		Model:   "qwen-plus",
		BaseURL: "https://dashscope.aliyuncs.com/compatible-mode/v1",
	})
	if err != nil {
		log.Fatalf("failed to create chat model: %v", err)
	}
	//示例1. 使用翻译模板
	translator := templates.Translator("en", "zh")
	mssages, _ := translator.Format(ctx, map[string]any{
		"text": "Hello, world!",
	})

	response, _ := chatModel.Generate(ctx, mssages)
	fmt.Println("翻译结果：", response.Content)

	//示例2. 使用代码审查模板
	codeReview := templates.CodeReview("go")
	mssages, _ = codeReview.Format(ctx, map[string]any{
		"language": "go",
		"code":     "package main\nimport \"fmt\"\nfunc main() {\n\tfmt.Println(\"Hello, world!\")\n}",
	})
	response, _ = chatModel.Generate(ctx, mssages)
	fmt.Println("代码审查结果：", response.Content)

}
