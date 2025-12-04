package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
	"log"
)

type UserProfile struct {
	Name      string
	Age       int
	Interests []string
	VIPLevel  int
}

func main() {
	ctx := context.Background()
	//
	////1.创建ChatTemplate
	//template := prompt.FromMessages(
	//	schema.FString,
	//	schema.SystemMessage("你是一个{role}"),
	//	schema.UserMessage("{prompt}"),
	//)
	////2.准备参数
	//variables := map[string]any{
	//	"role":   "AI助手",
	//	"prompt": "介绍一下GO编程语言",
	//}
	////3.调用ChatTemplate
	//message, err := template.Format(ctx, variables)
	//if err != nil {
	//	log.Fatalf("failed to format chat template: %v", err)
	//}
	////4.打印结果
	//for i, msg := range message {
	//	fmt.Printf("Message %d: %s\n", i, msg.Content)
	//}
	////5.使用生成的消息调用模型
	//chatModel, err := qwen.NewChatModel(ctx, &qwen.ChatModelConfig{
	//	APIKey:  "sk-525208983ed54ff993fcef91bd03a88a",
	//	Model:   "qwen-plus",
	//	BaseURL: "https://dashscope.aliyuncs.com/compatible-mode/v1",
	//})
	////6.获得响应
	//response, err := chatModel.Generate(ctx, message)
	//if err != nil {
	//	log.Fatalf("failed to generate response: %v", err)
	//}
	//fmt.Println("Response:", response.Content)
	//template := prompt.FromMessages(
	//	schema.FString,
	//	schema.SystemMessage("你是一个{role},你的特长是{ext}"),
	//	schema.UserMessage("{query}"),
	//	schema.AssistantMessage("我理解了，让我思考一下...", nil),
	//	schema.UserMessage("请详细说明"),
	//)
	//variables := map[string]any{
	//	"role":  "AI助手",
	//	"ext":   "回答各种问题",
	//	"query": "介绍一下GO编程语言",
	//}
	//message, err := template.Format(ctx, variables)
	//if err != nil {
	//	log.Fatalf("failed to format chat template: %v", err)
	//}
	//for i, msg := range message {
	//	fmt.Printf("Message %d: %s\n", i, msg.Content)
	//}
	template := prompt.FromMessages(
		schema.FString,
		schema.SystemMessage("你是一个{role},你的特长是{ext}"),
		schema.UserMessage(`用户信息：
			姓名：{name} 
			年龄：{age}
			兴趣：{interests} 
			VIP等级：{vip_level}
			请根据以上信息推荐合适内容`),
	)
	//准备用户数据
	user := UserProfile{
		Name:      "张三",
		Age:       30,
		Interests: []string{"编程", "旅行", "音乐"},
		VIPLevel:  5,
	}
	variables := map[string]any{
		"role":      "AI助手",
		"ext":       "回答各种问题",
		"name":      user.Name,
		"age":       user.Age,
		"interests": user.Interests,
		"vip_level": user.VIPLevel,
	}
	message, err := template.Format(ctx, variables)
	if err != nil {
		log.Fatalf("failed to format chat template: %v", err)
	}
	for i, msg := range message {
		fmt.Printf("Message %d: %s\n", i, msg.Content)
	}
}
