package cmd

import (
	"context"
	"fmt"
	"log"
	"nvoke/nvoke"
	"nvoke/pkg/embedding"

	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SearchSimilarEmbeddings(query string) {
	ctx := context.Background()

	if query == "" {
		log.Fatalf("Invalid query string \"%s\"\n", query)
	}
	// Initialize the OpenAI generator and vectorize the query
	client := openai.NewClient(OpenAIAPIKey)
	generator := embedding.NewOpenAIGenerator(client, openai.SmallEmbedding3, 1536)
	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI(MongoDBConnectionString)
	mongoClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(ctx)

	service := nvoke.NewRetrievalService(mongoClient, generator, client)
	items, err := service.SemanticSearch(ctx, nvoke.Query{Query: query, Persona: persona})
	if err != nil {
		log.Fatalf("Failed to find similar content %v", err)
	}
	fmt.Println("Similar documents")
	for _, item := range items {
		fmt.Printf("%v\n", item)
	}
}

var similarCmd = &cobra.Command{
	Use:   "similar",
	Short: "similarity search for text",
	Run: func(cmd *cobra.Command, args []string) {
		SearchSimilarEmbeddings(query)
	},
}

func init() {
	similarCmd.Flags().StringVarP(&query, "query", "q", "", "Text query to search similar embeddings")
	similarCmd.Flags().StringVarP(&persona, "persona", "p", "", "The persona to use when searching")
	similarCmd.Flags().IntVarP(&limit, "limit", "l", 10, "Max similar vectors limit.")
	similarCmd.Flags().IntVarP(&candidates, "candidates", "c", 200, "Number of candidates to consider.")
	rootCmd.AddCommand(similarCmd)
}
