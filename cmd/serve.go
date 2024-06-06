package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"nvoke/pkg/bible"
	"nvoke/pkg/embedding"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Setup an API to serve completion requests",
	Run: func(cmd *cobra.Command, args []string) {
		serve()
	},
}

var port int
var address string

func init() {
	serveCmd.Flags().IntVarP(&limit, "limit", "l", 10, "Max similar vectors limit.")
	serveCmd.Flags().IntVarP(&candidates, "candidates", "c", 200, "Number of candidates to consider.")
	serveCmd.Flags().IntVarP(&port, "port", "p", 80, "Listener port")
	serveCmd.Flags().StringVarP(&address, "address", "a", "0.0.0.0", "Listener address")

	rootCmd.AddCommand(serveCmd)
}

func serve() {
	r := chi.NewRouter()

	// c := cors.New(cors.Options{
	// 	AllowedOrigins: []string{"http://frontend.local"},
	// })

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/v1/completion", handleRequest)

	log.Printf("Listening on %s:%d", address, port)
	http.ListenAndServe(fmt.Sprintf("%s:%d", address, port), r)
}

type Query struct {
	Query string `json:"query"`
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var data Query
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Printf("Failed decode request body: %v\n", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if data.Query == "" {
		log.Println("data.Query is empty")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	client := openai.NewClient(OpenAIAPIKey)
	generator := embedding.NewOpenAIGenerator(client, openai.SmallEmbedding3, 1536)

	queryEmbedding, err := generator.GenerateEmbedding(ctx, data.Query)
	if err != nil {
		log.Printf("Failed to generate embedding for the query: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	clientOptions := options.Client().ApplyURI(MongoDBConnectionString)
	mongodb, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Printf("Failed to connect to MongoDB: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer mongodb.Disconnect(ctx)

	collection := mongodb.Database("bible").Collection("verses")
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
		log.Printf("Failed to find similar verses: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	req := openai.ChatCompletionRequest{
		Model: openai.GPT4o,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are Jesus. You will respond in language like that of the NKJV bible as if you are Jesus talking to his son.",
			},
		},
	}

	contextString := "Using the following verses for context to answer the question. Do not use other information or sources. \n context: "
	for cursor.Next(ctx) {
		var verse bible.Verse
		if err = cursor.Decode(&verse); err != nil {
			log.Printf("Failed to decode verse: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		contextString += fmt.Sprintf("%v %v:%v -- %v ", verse.Book, verse.Chapter, verse.Verse, verse.Text)
		req.Messages = append(req.Messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: fmt.Sprintf("%s \n question: %s", contextString, data.Query),
		})
	}

	response, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		log.Printf("Error generating completion: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	respData := map[string]interface{}{
		"answer": response.Choices[0].Message.Content,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respData)
}
