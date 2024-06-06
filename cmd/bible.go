package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"nvoke/pkg/bible"
	"nvoke/pkg/embedding"
	"os"
	"time"

	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

func GenerateAndSaveBibleEmbeddings() {
	// Load books from files (reusing the existing logic which might need to be refactored out from main)
	books := loadBooksFromFiles()

	// Initialize OpenAI client and necessary components
	client := openai.NewClient(OpenAIAPIKey)
	generator := embedding.NewOpenAIGenerator(client, openai.SmallEmbedding3, 1536)
	limiter := embedding.NewSteadyRateLimiter(2500, time.Minute, 10)
	adapter := bible.NewVerseAdapter()

	// Create and use the embedding service
	service := embedding.NewEmbeddingService(generator, adapter, limiter)

	verses := generateVerseDocuments(books)
	err := service.GenerateEmbeddings(context.Background(), verses)
	if err != nil {
		fmt.Printf("Failed to generate embeddings: %v\n", err)
		return
	}

	// Save verses with embeddings to JSON
	bytes, err := json.Marshal(verses)
	if err != nil {
		fmt.Printf("Failed to marshal verses: %v\n", err)
		return
	}
	file, err := os.Create("texts/bible/nkjv-verses.json")
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

func loadBooksFromFiles() []*bible.Book {
	books := make([]*bible.Book, 0)
	directory := "texts/bible/nkjv/"
	entries, err := os.ReadDir(directory)
	if err != nil {
		fmt.Println(err)
		return books
	}
	for _, entry := range entries {
		f, err := os.Open(directory + entry.Name())
		if err != nil {
			fmt.Printf("Error opening file %v\n", err)
			continue
		}
		defer f.Close()
		books = append(books, bible.ParseBooksFile(f)...)
	}
	return books
}
func generateVerseDocuments(books []*bible.Book) []*bible.Verse {
	content := make([]*bible.Verse, 0)
	for _, book := range books {
		for _, chapter := range book.Chapters {
			content = append(content, chapter.Verses...)
		}
	}
	return content
}

var generateBibleCmd = &cobra.Command{
	Use:   "bible",
	Short: "Generate embeddings for bible verses",
	Run: func(cmd *cobra.Command, args []string) {
		// Assuming you have a function `GenerateAndSaveEmbeddings` implemented
		GenerateAndSaveBibleEmbeddings()
	},
}

func init() {
	generateCmd.AddCommand(generateBibleCmd)
}
