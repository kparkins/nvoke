package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"nvoke/nvoke"
	"nvoke/pkg/embedding"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
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

// TODO
// upload.go

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
	ctx := context.Background()
	openaiClient := openai.NewClient(OpenAIAPIKey)
	generator := embedding.NewOpenAIGenerator(openaiClient, openai.SmallEmbedding3, 1536)

	clientOptions := options.Client().ApplyURI(MongoDBConnectionString)
	mongodb, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Printf("Failed to connect to MongoDB: %v\n", err)
		return
	}
	defer mongodb.Disconnect(ctx)
	service := nvoke.NewRetrievalService(mongodb, generator, openaiClient)

	// c := cors.New(cors.Options{
	// 	AllowedOrigins: []string{"http://frontend.local"},
	// })

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/v1/completion", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var data nvoke.Query
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			log.Printf("Failed decode request body: %v\n", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		completion, err := service.CreateChatCompletion(ctx, data)
		switch err {
		case nvoke.ErrInvalidQueryParameters:
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		case nvoke.ErrSimilaritySearchFailed,
			nvoke.ErrEmbeddingGenerationFailed,
			nvoke.ErrChatCompletionContextBuildFailed,
			nvoke.ErrChatCompletionFailed:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"answer": completion,
		})
	})

	r.Post("/v1/completion/stream", func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Failed to upgrade to WebSocket:", err)
			return
		}
		defer conn.Close()

		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Failed to read message from WebSocket:", err)
			return
		}

		var query nvoke.Query
		err = json.Unmarshal(message, &query)
		if err != nil {
			log.Println("Failed to unmarshal JSON message:", err)
		}

		stream, err := service.CreateChatCompletionStream(ctx, query)
		if err != nil {
			log.Println("Failed to create chat completion stream:", err)
			return
		}
		defer stream.Close()

		for {
			response, err := stream.Recv()
			if err != nil {
				log.Println("Failed to receive message from OpenAI:", err)
				return
			}

			err = conn.WriteMessage(websocket.TextMessage, []byte(response.Choices[0].Delta.Content))
			if err != nil {
				log.Println("Failed to write message to WebSocket:", err)
				return
			}
		}

	})

	log.Printf("Listening on %s:%d", address, port)
	http.ListenAndServe(fmt.Sprintf("%s:%d", address, port), r)
}
