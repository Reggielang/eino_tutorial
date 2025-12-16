package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudwego/eino-ext/components/model/qwen"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
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

	//2.创建工具
	//获取当前时间的工具
	searchTool := utils.NewTool(&schema.ToolInfo{
		Name: "search",
		Desc: "搜索信息",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"query": {
				Type:     "string",
				Desc:     "需要搜索的内容",
				Required: true,
			},
		}),
	},
		func(ctx context.Context, params map[string]any) (string, error) {
			query := params["query"].(string)
			result := fmt.Sprintf("找到关于'%s'的信息： Go是Google开发的编程语言", query)
			return result, nil
		})

	//4.创建React Agent
	reactAgent, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: chatModel,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: []tool.BaseTool{searchTool},
		},
	})
	if err != nil {
		log.Fatalf("failed to create react agent: %v", err)
	}

	//4. 使用agent
	messages := []*schema.Message{
		schema.SystemMessage("请告诉我GO 语言的特点"),
	}
	fmt.Println("[用户输入] 请告诉我GO 语言的特点")

	streamOutput, err := reactAgent.Stream(ctx, messages)
	if err != nil {
		log.Fatalf("failed to run react agent: %v", err)
	}
	defer streamOutput.Close()

	// 读取流式输出
	for {
		chunk, err := streamOutput.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Fatalf("failed to read stream output: %v", err)
		}
		fmt.Println(chunk.Content)
	}
}
