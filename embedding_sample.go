package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"log"
	"math"
	"os"
)

func main() {
	ctx := context.Background()

	//创建向量模型
	embedder, err := ark.NewEmbedder(ctx, &ark.EmbeddingConfig{
		BaseURL: "https://ark.cn-beijing.volces.com/api/v3/embeddings/multimodal",
		APIKey:  os.Getenv("ARK_KEY"),
		Model:   os.Getenv("ARK_EMBEDDING_MODEL"),
	})
	if err != nil {
		log.Fatalf("failed to create embedder: %v", err)
	}

	//向量的文本
	texts := []string{
		"go is a programming language",
		"python is a programming language",
		"今天天气很好",
	}
	//获取向量
	vectors, err := embedder.EmbedStrings(ctx, texts)
	if err != nil {
		log.Fatalf("failed to embed strings: %v", err)
	}
	//打印向量
	for i, text := range texts {
		fmt.Printf("text: %s, vector: %v\n", text, vectors[i])
	}

	//计算相似度
	similarity1 := cosineSimilarity(vectors[0], vectors[1])
	similarity2 := cosineSimilarity(vectors[0], vectors[2])
	fmt.Printf("similarity1: %f, similarity2: %f\n", similarity1, similarity2)

}

// 计算余弦相似度
func cosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) {
		return 0
	}
	var dotProduct, normA, normB float64
	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}
