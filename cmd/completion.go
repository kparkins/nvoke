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

var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Perform a similarity search and generate a completion for the most similar text",
	Run: func(cmd *cobra.Command, args []string) {
		SearchSimilarVersesAndGenerateCompletion(query)
	},
}

var completionContext string

func init() {
	completionCmd.Flags().StringVarP(&query, "query", "q", "", "Text query to search similar Bible verses")
	completionCmd.Flags().StringVarP(&completionContext, "context", "x", "", "Text context for the completion")
	completionCmd.Flags().IntVarP(&limit, "limit", "l", 10, "Max similar vectors limit.")
	completionCmd.Flags().IntVarP(&candidates, "candidates", "c", 200, "Number of candidates to consider.")
	rootCmd.AddCommand(completionCmd)
}

func SearchSimilarVersesAndGenerateCompletion(query string) {
	ctx := context.Background()

	if query == "" {
		log.Fatalf("Invalid query string \"%s\"\n", query)
	}

	client := openai.NewClient(OpenAIAPIKey)
	generator := embedding.NewOpenAIGenerator(client, openai.SmallEmbedding3, 1536)
	var err error
	var queryEmbedding []float32
	if completionContext != "" {
		queryEmbedding, err = generator.GenerateEmbedding(ctx, completionContext)
	} else {
		queryEmbedding, err = generator.GenerateEmbedding(ctx, query)
	}
	if err != nil {
		log.Fatalf("Failed to generate embedding for the query: %v", err)
	}

	clientOptions := options.Client().ApplyURI(MongoDBConnectionString)
	mongoClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(ctx)

	collection := mongoClient.Database("biblical").Collection("verses")
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

	cursor, err := collection.Aggregate(ctx, filter)
	if err != nil {
		log.Fatalf("Failed to find similar verses: %v", err)
	}
	defer cursor.Close(ctx)

	req := openai.ChatCompletionRequest{
		Model: openai.GPT4Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleSystem,
				//Content: "You are a biblical scholar who uses bible verses and references to them to answer questions. In particular you know the NKJV bible.",
				Content: "You are Jesus. You will respond in language like that of the NKJV bible as if you are Jesus talking to his son.",
			},
		},
	}
	contextString := "Using the following verses for context, answer the question. \n context: "
	for cursor.Next(ctx) {
		var verse bible.Verse
		if err = cursor.Decode(&verse); err != nil {
			log.Fatalf("Failed to decode verse: %v", err)
		}
		contextString += fmt.Sprintf("%v %v:%v -- %v ", verse.Book, verse.Chapter, verse.Verse, verse.Text)
	}

	req.Messages = append(req.Messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: fmt.Sprintf("%s \n question: %s", contextString, query),
	})

	response, err := client.CreateChatCompletion(ctx, req)

	if err != nil {
		log.Printf("error generating completion %v", err)
	}

	fmt.Printf("Answer: %s", response.Choices[0].Message.Content)
}
