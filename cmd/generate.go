package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"nvoke/pkg/bible"
	"nvoke/pkg/embedding"
	"nvoke/pkg/tao"
	"os"
	"time"

	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

func GenerateAndSaveEmbeddings[T any](adapter embedding.Adapter[T], content []T, writer io.Writer) {
	// Initialize OpenAI client and necessary components
	client := openai.NewClient(OpenAIAPIKey)
	generator := embedding.NewOpenAIGenerator(client, openai.SmallEmbedding3, 1536)
	limiter := embedding.NewSteadyRateLimiter(2500, time.Minute, 10)

	// Create and use the embedding embedder
	embedder := embedding.NewService(generator, adapter, limiter)

	err := embedder.GenerateEmbeddings(context.Background(), content)
	if err != nil {
		fmt.Printf("Failed to generate embeddings: %v\n", err)
		return
	}

	// Save verses with embeddings to JSON
	bytes, err := json.Marshal(content)
	if err != nil {
		fmt.Printf("Failed to marshal chapters: %v\n", err)
		return
	}

	_, err = writer.Write(bytes)
	if err != nil {
		fmt.Printf("Failed to write to chapters.json: %v\n", err)
	}
	fmt.Println("Embeddings generated and saved to chapters.json")
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate embeddings for text",
	Run: func(cmd *cobra.Command, args []string) {
		switch persona {
		case "bible":
			parser := bible.Parser{}
			verses := parser.Parse("")
			file, err := os.Create("texts/bible/nkjv-verses.json")
			if err != nil {
				fmt.Printf("Failed to create nkjv-verses.json: %v\n", err)
				return
			}
			defer file.Close()
			GenerateAndSaveEmbeddings(bible.NewEmbeddingAdapter(), verses, file)
		case "tao":
			parser := tao.Parser{}
			chapters := parser.Parse("")
			file, err := os.Create("texts/tao/linn/chapters.json")
			if err != nil {
				fmt.Printf("Failed to create chapters.json: %v\n", err)
				return
			}
			defer file.Close()
			GenerateAndSaveEmbeddings(tao.NewEmbeddingAdapter(), chapters, file)
		}
	},
}

func init() {
	similarCmd.Flags().StringVarP(&persona, "persona", "p", "", "The persona to use when searching")
	rootCmd.AddCommand(generateCmd)
}
