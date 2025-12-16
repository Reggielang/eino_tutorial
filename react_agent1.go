package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino-ext/components/model/qwen"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
	"log"
	"time"
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
	getTimeTool := utils.NewTool(&schema.ToolInfo{
		Name:        "get_time",
		Desc:        "获取当前时间",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{}),
	},
		func(ctx context.Context, params map[string]any) (string, error) {
			now := time.Now().Format("2006-01-02 15:04:05")
			fmt.Printf("[工具执行] get_time ->%s\n", now)
			return now, nil
		})
	//3.简单计算器工具
	calculatorTool := utils.NewTool(&schema.ToolInfo{
		Name: "calculator",
		Desc: "计算器工具",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"expression": {
				Type:     "string",
				Desc:     "需要计算的表达式",
				Required: true,
			},
		}),
	},
		func(ctx context.Context, params map[string]any) (string, error) {
			expr := params["expression"].(string)
			//简化：这里只演示，实际应该用真正的表达式解析
			result := "15"
			fmt.Printf("[工具执行] calculator ->%s\n", expr)
			return result, nil

		})
	//4.创建React Agent
	reactAgent, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: chatModel,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: []tool.BaseTool{getTimeTool, calculatorTool},
		},
	})
	if err != nil {
		log.Fatalf("failed to create react agent: %v", err)
	}

	//4. 使用agent
	messages := []*schema.Message{
		schema.SystemMessage("现在几点了?"),
	}
	fmt.Println("[用户输入] 现在几点了?")

	output, err := reactAgent.Generate(ctx, messages)
	if err != nil {
		log.Fatalf("failed to run react agent: %v", err)
	}
	fmt.Println("回答:", output.Content)
}
