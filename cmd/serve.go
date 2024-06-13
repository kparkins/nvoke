package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"nvoke/pkg/embedding"
	"nvoke/pkg/nvoke"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
// generate.go
// completion.go
// similar.go
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
	service := nvoke.NewCompletionService(mongodb, generator, openaiClient)

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
	},
	)

	log.Printf("Listening on %s:%d", address, port)
	http.ListenAndServe(fmt.Sprintf("%s:%d", address, port), r)
}
