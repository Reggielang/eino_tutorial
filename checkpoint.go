package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino-ext/components/model/qwen"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"log"
	"sync"
)

var _ compose.CheckPointStore = (*memoryCheckPointStore)(nil)

// 简单的内存 checkpoint 存储
type memoryCheckPointStore struct {
	mu   sync.RWMutex
	data map[string][]byte
}

func newMemoryCheckPointStore() *memoryCheckPointStore {
	return &memoryCheckPointStore{
		data: make(map[string][]byte),
	}
}

func (s *memoryCheckPointStore) Get(ctx context.Context, checkPointID string) ([]byte, bool, error) {
	s.mu.RLock()
	data, ok := s.data[checkPointID]
	return data, ok, nil
}
func (s *memoryCheckPointStore) Set(ctx context.Context, checkPointID string, data []byte) error {
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

	//2. 创建 Agent
	agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "BookRecommender",
		Description: "书籍推荐 agent",
		Instruction: "你是一个书籍推荐专家，请根据用户需求推荐合适的书籍",
		Model:       chatModel,
	})
	if err != nil {
		log.Fatalf("failed to create techAgent: %v", err)
	}
	//3. 创建 checkpoint 存储
	store := newMemoryCheckPointStore()

	//4. 创建 runner
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           agent,
		EnableStreaming: false,
		CheckPointStore: store,
	})
	//5. 运行(带checkpoint ID)
	checkPointID := "session_001"
	query := "我想学习go语言，请推荐一些书籍"
	fmt.Printf("Query: %s\n", query)
	iter := runner.Run(ctx, []adk.Message{
		schema.UserMessage(query),
	}, adk.WithCheckPointID(checkPointID))
	//6. 处理结果
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
				fmt.Printf("[%s]: %s\n", event.AgentName, msg.Content)
			}
		}
	}
	////7. 重新运行
	//fmt.Println("重新运行")
	//resumeIter,err := runner.Resume(ctx, checkPointID)
	//if err != nil {
	//	log.Fatalf("failed to resume agent: %v", err)
	//}
	////处理恢复事件
}
