package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudwego/eino-ext/components/model/qwen"
	"github.com/cloudwego/eino/schema"
	"io"
	"log"
)

// QWEN_API_KEY=sk-525208983ed54ff993fcef91bd03a88a
func main() {
	//1.创建上下文
	ctx := context.Background()
	//2.创建chat model
	chatMode, err := qwen.NewChatModel(ctx, &qwen.ChatModelConfig{
		APIKey:  "sk-525208983ed54ff993fcef91bd03a88a",
		Model:   "qwen-plus",
		BaseURL: "https://dashscope.aliyuncs.com/compatible-mode/v1",
	})
	if err != nil {
		log.Fatalf("failed to create chat model: %v", err)
	}
	//3.准备message
	message := []*schema.Message{
		schema.SystemMessage("你是一个AI助手"),
		schema.UserMessage("请你介绍一下 GO 编程短文"),
	}
	//4.调用chat model - stream 流式生成
	stream, err := chatMode.Stream(ctx, message)
	if err != nil {
		log.Fatalf("failed to call chat model: %v", err)
	}
	//5.关闭流式
	defer stream.Close()

	fmt.Println("AI 回复：")
	//6. 逐块的接收并打印
	for {
		chunk, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				//流式结束
				break
			}
			log.Fatalf("failed to receive chunk: %v", err)
		}
		//打印内容
		fmt.Print(chunk.Content)
	}
	fmt.Println("\n 回复结束")

	////6. 逐块的接收并打印
	//var fullContent strings.Builder
	//fmt.Println("AI 回复：")
	//
	//for {
	//	chunk, err := stream.Recv()
	//	if err != nil {
	//		if errors.Is(err, io.EOF) {
	//			//流式结束
	//			break
	//		}
	//		log.Fatalf("failed to receive chunk: %v", err)
	//	}
	//	//打印内容
	//	fmt.Print(chunk.Content)
	//	//同时收集完整内容
	//	fullContent.WriteString(chunk.Content)
	//}
	//fmt.Println("\n 回复结束")
	//fmt.Println("完整内容：", fullContent.String())
}
