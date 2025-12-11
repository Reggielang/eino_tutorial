package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

type Calculator struct{}

// Info 返回工具信息
func (c *Calculator) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "Calculator",
		Desc: "执行基本数学计算（加、减、乘、除）",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"operation": {
				Type:     "string",
				Desc:     "运行类型： add, sub, mul, div",
				Required: true,
			},
			"num1": {
				Type:     "number",
				Desc:     "第一个数字",
				Required: true,
			},
			"num2": {
				Type:     "number",
				Desc:     "第二个数字",
				Required: true,
			},
		}),
	}, nil
}

// 参数结构
type CalculatorRequest struct {
	Operation string  `json:"operation"`
	Num1      float64 `json:"num1"`
	Num2      float64 `json:"num2"`
}

// 结果结构
type CalculatorResponse struct {
	Result float64 `json:"result"`
	Error  string  `json:"error,omitempty"`
}

// InvokableRun 执行计算
func (t *Calculator) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	//1.解析参数
	var request CalculatorRequest
	if err := json.Unmarshal([]byte(argumentsInJSON), &request); err != nil {
		return "", fmt.Errorf("failed to unmarshal arguments: %v", err)
	}

	//2.执行计算
	var result float64
	switch request.Operation {
	case "add":
		result = request.Num1 + request.Num2
	case "sub":
		result = request.Num1 - request.Num2
	case "mul":
		result = request.Num1 * request.Num2
	case "div":
		if request.Num2 == 0 {
			return "", fmt.Errorf("division by zero")
		}
		result = request.Num1 / request.Num2

	default:
		resultJSON, _ := json.Marshal(CalculatorResponse{Error: "invalid operation"})
		return string(resultJSON), nil
	}
	//3.返回结果
	resultJSON, err := json.Marshal(CalculatorResponse{Result: result})
	if err != nil {
		return "", err
	}
	return string(resultJSON), nil
}

func main() {
	ctx := context.Background()
	calculator := &Calculator{}

	//测试工具
	testCases := []struct {
		operation  string
		num1, num2 float64
	}{
		{"add", 1, 2},
		{"sub", 5, 3},
		{"mul", 4, 6},
		{"div", 8, 2},
	}

	for _, testCase := range testCases {
		request := CalculatorRequest{
			Operation: testCase.operation,
			Num1:      testCase.num1,
			Num2:      testCase.num2,
		}

		argumentsJSON, _ := json.Marshal(request)
		result, err := calculator.InvokableRun(ctx, string(argumentsJSON))
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		fmt.Printf("Result: %s\n", result)
	}
}
