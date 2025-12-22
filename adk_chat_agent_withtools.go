package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino-ext/components/model/qwen"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
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
	//2. 创建工具
	timeTool := utils.NewTool(&schema.ToolInfo{
		Name:        "time",
		Desc:        "Get current time",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{}),
	},
		func(ctx context.Context, params map[string]any) (string, error) {
			return time.Now().Format("2006-01-02 15:04:05"), nil

		})
	// 计算器工具
	calcTool := utils.NewTool(&schema.ToolInfo{
		Name: "calc",
		Desc: "Calculate the result of the given expression",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"expression": {
				Type:     schema.String,
				Desc:     "The expression to calculate",
				Required: true,
			},
		}),
	}, func(ctx context.Context, params map[string]any) (string, error) {
		_, ok := params["expression"].(string)
		if !ok {
			return "", fmt.Errorf("invalid expression")
		}
		result := eval()
		return fmt.Sprintf("%f", result), nil

	},
	)
	//2. 创建ChatModelAgent
	chatAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "ToolAssistant",
		Description: "一个能够使用工具的助手Agent",
		Instruction: `你是一个智能助手，可以使用以下工具:
						- time: 获取当前时间
						- calc: 执行数学计算
当用户需要获取时间或者进行计算时，请调用相应的工具`,
		Model: chatModel,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []tool.BaseTool{timeTool, calcTool},
			},
		},
		MaxIterations: 10,
	})
	if err != nil {
		log.Fatalf("failed to create chat model agent: %v", err)
	}
	//3. 创建 Runner
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           chatAgent,
		EnableStreaming: false,
	})
	//4. 运行Agent
	queries := []string{"现在几点了?", "帮我计算123+456"}

	for _, query := range queries {
		fmt.Printf("Query: %s\n", query)
		iter := runner.Query(ctx, query)
		for {
			event, ok := iter.Next()
			if !ok {
				//迭代器关闭，退出循环
				break
			}
			if event.Err != nil {
				log.Fatalf("failed to run agent: %v", event.Err)
			}
			if event.Output != nil && event.Output.MessageOutput != nil {
				msg := event.Output.MessageOutput.Message
				if msg != nil {
					// 检查是否有工具调用
					if len(msg.ToolCalls) > 0 {
						for _, toolCall := range msg.ToolCalls {
							// 检查是否有工具调用
							fmt.Printf("Tool Call: %s\n", toolCall.Function.Name)
						}
					} else if msg.Content != "" {
						fmt.Printf("Response: %s\n", msg.Content)
					}
				}
			}
		}
	}

}

func eval() float64 {
	return 0.0
}
