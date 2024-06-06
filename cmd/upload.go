package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"nvoke/pkg/bible"
	"nvoke/pkg/tao"
	"os"

	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func UploadBibleVersesToMongoDB() {
	// Open the JSON file
	file, err := os.Open("texts/bible/nkjv-verses.json")
	if err != nil {
		fmt.Printf("Failed to open verses.json: %v\n", err)
		return
	}
	defer file.Close()

	// Read and decode the file content
	var verses []*bible.Verse
	err = json.NewDecoder(file).Decode(&verses)
	if err != nil {
		fmt.Printf("Failed to decode verses: %v\n", err)
		return
	}

	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI(MongoDBConnectionString)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		fmt.Printf("Failed to connect to MongoDB: %v\n", err)
		return
	}
	defer client.Disconnect(context.Background())

	// Insert documents
	collection := client.Database("bible").Collection("verses")
	documents := make([]interface{}, len(verses))
	for i, verse := range verses {
		documents[i] = verse
	}
	_, err = collection.InsertMany(context.Background(), documents)
	if err != nil {
		fmt.Printf("Failed to insert verses into MongoDB: %v\n", err)
		return
	}

	fmt.Println("Verses uploaded to MongoDB successfully")
}

func UploadTaoChaptersToMongoDB() {
	// Open the JSON file
	// file, err := os.Open("texts/tao/linn/tao.json")
	// if err != nil {
	// 	fmt.Printf("Failed to open verses.json: %v\n", err)
	// 	return
	// }
	// defer file.Close()

	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI(MongoDBConnectionString)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		fmt.Printf("Failed to connect to MongoDB: %v\n", err)
		return
	}
	defer client.Disconnect(context.Background())

	// Insert documents
	collection := client.Database("tao").Collection("chapters")

	document := tao.Chapter{
		Chapter: 1,
		Text: `The Tao that can be spoken is not the eternal Tao
		The name that can be named is not the eternal name
		The nameless is the origin of Heaven and Earth
		The named is the mother of myriad things
		Thus, constantly without desire, one observes its essence
		Constantly with desire, one observes its manifestations
		These two emerge together but differ in name
		The unity is said to be the mystery
		Mystery of mysteries, the door to all wonders`,
		Embedding: nil,
	}
	_, err = collection.InsertOne(context.Background(), document)
	if err != nil {
		fmt.Printf("Failed to insert chapters into MongoDB: %v\n", err)
		return
	}

	fmt.Println("Chapters uploaded to MongoDB successfully")
}

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload documents to MongoDB",
	Run: func(cmd *cobra.Command, args []string) {
		// Assuming you have a function `UploadVersesToMongoDB` implemented
		UploadBibleVersesToMongoDB()
		UploadTaoChaptersToMongoDB()
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
}
