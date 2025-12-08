package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino/compose"
	"log"
	"strings"
)

func main() {
	ctx := context.Background()

	//创建一个简单的Chain: 输入字符串=转大写=添加前缀=输出
	chain := compose.NewChain[string, string]()
	chain.
		//lambda1： 转大写
		AppendLambda(compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
			return strings.ToUpper(input), nil
		})).
		//lambda2： 添加前缀
		AppendLambda(compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
			return "Prefix: " + input, nil
		}))

	runnable, err := chain.Compile(ctx)
	if err != nil {
		log.Fatalf("failed to compile chain: %v", err)
	}

	output, err := runnable.Invoke(ctx, "hello world")
	if err != nil {
		log.Fatalf("failed to invoke chain: %v", err)
	}
	fmt.Println("chain result:", output)
}
