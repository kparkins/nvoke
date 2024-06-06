package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"nvoke/pkg/embedding"
	"nvoke/pkg/tao"
	"os"
	"time"

	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

func GenerateAndSaveTaoEmbeddings() {
	// Load books from files (reusing the existing logic which might need to be refactored out from main)
	chapters := loadChaptersFromFile()

	// Initialize OpenAI client and necessary components
	client := openai.NewClient(OpenAIAPIKey)
	generator := embedding.NewOpenAIGenerator(client, openai.SmallEmbedding3, 1536)
	limiter := embedding.NewSteadyRateLimiter(2500, time.Minute, 10)
	adapter := tao.NewAdapter()

	// Create and use the embedding service
	service := embedding.NewEmbeddingService(generator, adapter, limiter)

	err := service.GenerateEmbeddings(context.Background(), chapters)
	if err != nil {
		fmt.Printf("Failed to generate embeddings: %v\n", err)
		return
	}

	// Save verses with embeddings to JSON
	bytes, err := json.Marshal(chapters)
	if err != nil {
		fmt.Printf("Failed to marshal verses: %v\n", err)
		return
	}
	file, err := os.Create("texts/tao/linn/chapters.json")
	if err != nil {
		fmt.Printf("Failed to create verses.json: %v\n", err)
		return
	}
	defer file.Close()

	_, err = file.Write(bytes)
	if err != nil {
		fmt.Printf("Failed to write to verses.json: %v\n", err)
	}
	fmt.Println("Embeddings generated and saved to verses.json")
}

func loadChaptersFromFile() []*tao.Chapter {
	chapters := make([]*tao.Chapter, 0)

	return chapters
}

var generateTaoCmd = &cobra.Command{
	Use:   "tao",
	Short: "Generate embeddings for Tao Te Ching chapters",
	Run: func(cmd *cobra.Command, args []string) {
		// Assuming you have a function `GenerateAndSaveEmbeddings` implemented
		GenerateAndSaveTaoEmbeddings()
	},
}

func init() {
	generateCmd.AddCommand(generateTaoCmd)
}
