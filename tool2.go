package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
	"log"
	"time"
)

type TimeParam struct {
	Format string `json:"format"`
}
type TimeRequest struct {
	CurrentTime string `json:"current_time"`
}

func GetCurrentTime(ctx context.Context, params *TimeParam) (*TimeRequest, error) {
	now := time.Now()
	var result string
	switch params.Format {
	case "date":
		result = now.Format("2006-01-02")
	case "time":
		result = now.Format("15:04:05")
	default:
		result = now.Format("2006-01-02 15:04:05")
	}
	return &TimeRequest{CurrentTime: result}, nil

}

func main() {
	// 测试工具
	ctx := context.Background()
	//使用NewTool()创建一个工具
	timeTool := utils.NewTool(&schema.ToolInfo{
		Name: "get_time",
		Desc: "获取当前时间",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"format": {
				Type:     schema.String,
				Desc:     "时间格式：date(日期),time(时间),datetime(日期时间)",
				Required: false,
			},
		}),
	}, GetCurrentTime)

	//测试工具
	testFormats := []string{"date", "time", "datetime"}
	for _, format := range testFormats {
		request := TimeParam{Format: format}
		b, _ := json.Marshal(request)
		//工具执行
		response, err := timeTool.InvokableRun(ctx, string(b))
		if err != nil {
			log.Printf("failed to run tool: %v", err)
		}
		fmt.Printf("Format: %s, Current Time: %s\n", format, response)
	}
}
