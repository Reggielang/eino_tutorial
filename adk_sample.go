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
	"sync"
)

var _ compose.CheckPointStore = (*memoryStore)(nil)

// 简单的内存 checkpoint 存储
type memoryStore struct {
	mu   sync.RWMutex
	data map[string][]byte
}

func newmemoryStore() *memoryStore {
	return &memoryStore{data: make(map[string][]byte)}
}

func (s *memoryStore) Get(ctx context.Context, checkPointID string) ([]byte, bool, error) {
	s.mu.RLock()
	data, ok := s.data[checkPointID]
	return data, ok, nil
}
func (s *memoryStore) Set(ctx context.Context, checkPointID string, data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[checkPointID] = data
	return nil
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

	//2. 创建 工具
	bookSearchTool := utils.NewTool(
		&schema.ToolInfo{
			Name: "bookSearch",
			Desc: "搜索书籍信息",
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
				"keyword": {
					Type:     schema.String,
					Desc:     "搜索关键词",
					Required: true,
				},
			}),
		},
		func(ctx context.Context, params map[string]any) (string, error) {
			keyword := params["keyword"].(string)
			//模拟书籍搜索
			books := map[string][]string{
				"Go":     {"《Go语言程序设计》", "《Go 并发编程实战》", "《Go 语言高级编程》"},
				"Python": {"《Python 编程：从入门到实践》", "《流畅的Python》", "《Python核心编程》"},
			}
			if results, ok := books[keyword]; ok {
				return fmt.Sprintf("搜索到以下书籍：%s", results), nil
			}
			return fmt.Sprintf("未找到关键词 '%s' 相关书籍", keyword), nil
		},
	)
	//3. 创建工具：询问澄清
	askTool := utils.NewTool(
		&schema.ToolInfo{
			Name: "ask_for_clarification",
			Desc: "询问用户以获取更多信息",
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
				"question": {
					Type:     schema.String,
					Desc:     "要询问用户的问题",
					Required: true,
				},
			}),
		},
		func(ctx context.Context, params map[string]any) (string, error) {
			question := params["question"].(string)
			return fmt.Sprintf("用户需要澄清: %s", question), nil
		})
	//4. 创建agent
	agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "BookRecommender",
		Description: "书籍推荐 agent,能够搜索书籍并询问用户偏好",
		Instruction: `你是一个书籍推荐专家，你可以：
						1. 使用 bookSearch 工具搜索书籍
						2. 使用 ask_for_clarification 工具询问用户的偏好
						3. 使用搜索结果和用户偏好推荐合适的书籍`,
		Model: chatModel,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []tool.BaseTool{bookSearchTool, askTool},
			},
		},
		Exit:          adk.ExitTool{}, //使用Exit工具
		MaxIterations: 10,
	})
	if err != nil {
		log.Fatalf("failed to create techAgent: %v", err)
	}
	//5. 创建 checkpoint 存储
	store := newmemoryStore()

	//6. 创建 runner
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           agent,
		EnableStreaming: false,
		CheckPointStore: store,
	})
	//7. 运行(带checkpoint ID)
	checkPointID := "session_001"
	query := "我想学习Go语言"
	fmt.Printf("Query: %s\n", query)
	iter := runner.Run(ctx, []adk.Message{
		schema.UserMessage(query),
	}, adk.WithCheckPointID(checkPointID))
	//8. 处理结果
	for {
		event, ok := iter.Next()
		if !ok {
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
	///9. 重新运行
	//fmt.Println("重新运行")
	//resumeIter,err := runner.Resume(ctx, checkPointID)
	//if err != nil {
	//	log.Fatalf("failed to resume agent: %v", err)
	//}
	////处理恢复事件
}
