package cmd

import (
	"context"
	"fmt"
	"log"
	"nvoke/pkg/embedding"
	"nvoke/pkg/nvoke"

	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ragCmd = &cobra.Command{
	Use:   "rag",
	Short: "Perform a similarity search and generate a completion for the most similar text",
	Run: func(cmd *cobra.Command, args []string) {
		RetrievalAugmentedSearch(query)
	},
}

var completionContext string

func init() {
	ragCmd.Flags().StringVarP(&query, "query", "q", "", "Text query to search similar Bible verses")
	ragCmd.Flags().StringVarP(&completionContext, "context", "x", "", "Text context for the completion")
	ragCmd.Flags().StringVarP(&persona, "persona", "p", "", "The persona to use when searching")
	ragCmd.Flags().IntVarP(&limit, "limit", "l", 10, "Max similar vectors limit.")
	ragCmd.Flags().IntVarP(&candidates, "candidates", "c", 200, "Number of candidates to consider.")
	rootCmd.AddCommand(ragCmd)
}

func RetrievalAugmentedSearch(query string) {
	ctx := context.Background()

	if query == "" {
		log.Fatalf("Invalid query string \"%s\"\n", query)
	}

	openaiClient := openai.NewClient(OpenAIAPIKey)
	generator := embedding.NewOpenAIGenerator(openaiClient, openai.SmallEmbedding3, 1536)

	data := nvoke.Query{
		Query:   query,
		Persona: persona,
	}
	if completionContext != "" {
		data.Query = completionContext
	}

	clientOptions := options.Client().ApplyURI(MongoDBConnectionString)
	mongodb, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongodb.Disconnect(ctx)

	service := nvoke.NewCompletionService(mongodb, generator, openaiClient)

	completion, err := service.CreateChatCompletion(ctx, data)
	if err != nil {
		log.Fatalf("Error creating completion %v\n", err)
	}
	fmt.Printf("Completion: %v\n", completion)
}
