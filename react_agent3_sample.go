package main

import (
	"context"
	"encoding/json"
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

// 餐厅数据结构
type Restaurant struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Cuisine  string   `json:"cuisine"`
	Location string   `json:"location"`
	Rating   float64  `json:"rating"`
	Tags     []string `json:"tags"`
}

type Dish struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Desc  string  `json:"desc"`
	Spicy bool    `json:"spicy"`
}

//模拟餐厅数据库

var restaurantDB = []Restaurant{
	{
		ID:       "1",
		Name:     "川菜馆",
		Cuisine:  "川菜",
		Location: "北京市朝阳区",
		Rating:   4.5,
		Tags:     []string{"川菜", "辣", "重口味"},
	},
	{
		ID:       "2",
		Name:     "粤菜馆",
		Cuisine:  "粤菜",
		Location: "上海市浦东新区",
		Rating:   4.8,
		Tags:     []string{"粤菜", "清淡", "美味"},
	},
	{
		ID:       "3",
		Name:     "湘菜馆",
		Cuisine:  "湘菜",
		Location: "广州市天河区",
		Rating:   4.3,
		Tags:     []string{"湘菜", "辣", "重口味"},
	},
}

var dishDB = map[string][]Dish{
	"1": {
		{
			Name:  "水煮鱼",
			Price: 68.0,
			Desc:  "川菜经典菜品，麻辣鲜香，口感鲜美。",
			Spicy: true,
		},
		{
			Name:  "宫保鸡丁",
			Price: 38.0,
			Desc:  "川菜经典菜品，口感鲜美，回味无穷。",
			Spicy: false,
		},
		{
			Name:  "麻婆豆腐",
			Price: 28.0,
			Desc:  "川菜经典菜品，麻辣鲜香，口感鲜美。",
			Spicy: true,
		},
	},
	"2": {
		{
			Name:  "白切鸡",
			Price: 88.0,
			Desc:  "粤菜经典菜品，肉质鲜嫩，清淡可口。",
			Spicy: false,
		},
		{
			Name:  "叉烧",
			Price: 58.0,
			Desc:  "粤菜经典菜品，肉质鲜嫩，口感鲜美。",
			Spicy: false,
		},
		{
			Name:  "煲仔饭",
			Price: 38.0,
			Desc:  "粤菜经典菜品，口感鲜美，回味无穷。",
			Spicy: false,
		},
	},
	"3": {
		{
			Name:  "剁椒鱼头",
			Price: 68.0,
			Desc:  "湘菜经典菜品，麻辣鲜香，口感鲜美。",
			Spicy: true,
		},
		{
			Name:  "口水鸡",
			Price: 48.0,
			Desc:  "湘菜经典菜品，口感鲜美，回味无穷。",
			Spicy: true,
		},
		{
			Name:  "辣子鸡",
			Price: 38.0,
			Desc:  "湘菜经典菜品，麻辣鲜香，口感鲜美。",
			Spicy: true,
		},
	},
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

	//2.创建工具
	//查询餐厅工具
	restaurantTool := utils.NewTool(&schema.ToolInfo{
		Name: "query_restaurant",
		Desc: "查询餐厅信息(根据位置，菜系，口味)",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"location": {
				Type: schema.String,
				Desc: "餐厅位置,例如：北京、上海",
			},
			"cuisine": {
				Type: schema.String,
				Desc: "餐厅菜系,例如：川菜、粤菜",
			},
			"spicy": {
				Type: schema.Boolean,
				Desc: "是否辣，例如：true、false",
			},
		}),
	},
		func(ctx context.Context, params map[string]any) (string, error) {
			fmt.Printf("[餐厅工具] 查询餐厅信息(根据位置，菜系，口味) 参数：%v\n", params)
			var results []Restaurant
			location, _ := params["location"].(string)
			cuisine, _ := params["cuisine"].(string)
			spicy, _ := params["spicy"].(bool)

			for _, r := range restaurantDB {
				match := true
				if location != "" && r.Location != location {
					match = false
				}
				if cuisine != "" && r.Cuisine != cuisine {
					match = false
				}
				if spicy {
					hasSpicy := false
					for _, tag := range r.Tags {
						if tag == "辣" {
							hasSpicy = true
							break
						}
					}
					if !hasSpicy {
						match = false
					}
				}

				if match {
					results = append(results, r)
				}
			}
			resultJSON, _ := json.Marshal(results)
			return string(resultJSON), nil
		})
	//查询菜品工具
	dishTool := utils.NewTool(&schema.ToolInfo{
		Name: "query_dish",
		Desc: "查询菜品信息(根据餐厅ID，菜品名称，菜品价格)",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"restaurant_id": {
				Type:     schema.String,
				Desc:     "餐厅ID,例如：1、2、3",
				Required: true,
			},
		}),
	},
		func(ctx context.Context, params map[string]any) (string, error) {
			fmt.Printf("[菜品工具] 查询菜品信息(根据餐厅ID，菜品名称，菜品价格) 参数：%v\n", params)
			restaurantID := params["restaurant_id"].(string)
			dishes, exists := dishDB[restaurantID]
			if !exists {
				return "", fmt.Errorf("restaurant not found")
			}
			resultJSON, _ := json.Marshal(dishes)
			return string(resultJSON), nil
		},
	)

	//4.创建React Agent
	reactAgent, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: chatModel,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: []tool.BaseTool{restaurantTool, dishTool},
		},
	})
	if err != nil {
		log.Fatalf("failed to create react agent: %v", err)
	}

	//4. 使用agent
	messages := []*schema.Message{
		schema.SystemMessage("我在北京，我想吃辣一点的，给我推荐餐厅和特色菜"),
	}
	fmt.Println("[用户输入] 我在北京，我想吃辣一点的，给我推荐餐厅和特色菜")

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
