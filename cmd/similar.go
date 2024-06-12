package cmd

import (
	"context"
	"fmt"
	"log"
	"nvoke/pkg/bible"
	"nvoke/pkg/embedding"

	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SearchSimilarVerses(query string) {
	ctx := context.Background()

	if query == "" {
		log.Fatalf("Invalid query string \"%s\"\n", query)
	}
	// Initialize the OpenAI generator and vectorize the query
	client := openai.NewClient(OpenAIAPIKey)
	generator := embedding.NewOpenAIGenerator(client, openai.SmallEmbedding3, 1536)
	queryEmbedding, err := generator.GenerateEmbedding(ctx, query)
	if err != nil {
		log.Fatalf("Failed to generate embedding for the query: %v", err)
	}

	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI(MongoDBConnectionString)
	mongoClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(ctx)

	// Prepare the MongoDB query for similarity search using cosine similarity
	collection := mongoClient.Database("bible").Collection("verses")
	filter := bson.A{
		bson.D{
			{Key: "$vectorSearch",
				Value: bson.D{
					{Key: "index", Value: "embedding"},
					{Key: "path", Value: "embedding"},
					{Key: "queryVector", Value: queryEmbedding},
					{Key: "numCandidates", Value: candidates},
					{Key: "limit", Value: limit},
				},
			},
		},
	}

	// Find similar verses
	cursor, err := collection.Aggregate(ctx, filter)
	if err != nil {
		log.Fatalf("Failed to find similar verses: %v", err)
	}
	defer cursor.Close(ctx)

	// Iterate through the results
	for cursor.Next(ctx) {
		var verse bible.Verse
		if err = cursor.Decode(&verse); err != nil {
			log.Fatalf("Failed to decode verse: %v", err)
		}
		fmt.Printf("Similar Verse: %v\n", verse.Text)
	}
	if err = cursor.Err(); err != nil {
		log.Fatalf("Error during cursor iteration: %v", err)
	}
}

var similarCmd = &cobra.Command{
	Use:   "similar",
	Short: "similarity search for bible verses",
	Run: func(cmd *cobra.Command, args []string) {
		SearchSimilarVerses(query)
	},
}

func init() {
	similarCmd.Flags().StringVarP(&query, "query", "q", "", "Text query to search similar embeddings")
	similarCmd.Flags().StringVarP(&persona, "persona", "p", "", "The persona to use when searching")
	similarCmd.Flags().IntVarP(&limit, "limit", "l", 10, "Max similar vectors limit.")
	similarCmd.Flags().IntVarP(&candidates, "candidates", "c", 200, "Number of candidates to consider.")
	rootCmd.AddCommand(similarCmd)
}
