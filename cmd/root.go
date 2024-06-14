package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var OpenAIAPIKey string
var MongoDBConnectionString string

var limit int
var candidates int
var query string
var persona string

var rootCmd = &cobra.Command{
	Use:   "nvoke",
	Short: "nvoke is a tool to manage embeddings and storage for a variety of religious and spiritual texts.",
}

// Execute executes the root command.
func Execute() {
	_ = godotenv.Load()

	// Get the environment variables
	OpenAIAPIKey = os.Getenv("OPENAI_API_KEY")
	MongoDBConnectionString = os.Getenv("MONGODB_CONNECTION_STRING_SRV")

	if OpenAIAPIKey == "" || MongoDBConnectionString == "" {
		log.Fatalf("Error loading environment variables")
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
