package nvoke

import (
	"context"
	"errors"
	"fmt"
	"log"
	"nvoke/pkg/embedding"

	"github.com/sashabaranov/go-openai"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var ErrInvalidQueryParameters = errors.New("invalid query data")
var ErrEmbeddingGenerationFailed = errors.New("embedding generation failed")
var ErrSimilaritySearchFailed = errors.New("similarity search failed")
var ErrChatCompletionContextBuildFailed = errors.New("failed to build completion context")
var ErrChatCompletionFailed = errors.New("failed to create chat completion")

// RetrievalService holds the parameters needed to serve completion requests.
type RetrievalService struct {
	Mongodb        *mongo.Client
	Generator      embedding.Generator
	OpenAI         *openai.Client
	EmbeddingModel string
	Limit          int
	Candidates     int
	KnowledgeBases map[string]KnowledgeBase
}

func NewRetrievalService(mongodb *mongo.Client, generator embedding.Generator, openai *openai.Client) *RetrievalService {
	return &RetrievalService{
		Mongodb:        mongodb,
		OpenAI:         openai,
		Generator:      generator,
		KnowledgeBases: KnowledgeBases,
	}
}

func (rs *RetrievalService) WithKnowledgeBases(knowledge map[string]KnowledgeBase) {
	rs.KnowledgeBases = knowledge
}

func (rs *RetrievalService) SemanticSearch(ctx context.Context, query Query) ([]interface{}, error) {
	if query.Query == "" {
		log.Println("data.Query is empty")
		return nil, ErrInvalidQueryParameters
	}
	knowledgeBase, ok := rs.KnowledgeBases[query.Persona]
	if !ok {
		log.Printf("invalid persona %v\n", query.Persona)
		return nil, ErrInvalidQueryParameters
	}

	queryEmbedding, err := rs.Generator.GenerateEmbedding(ctx, query.Query)
	if err != nil {
		log.Printf("Failed to generate embedding for the query: %v\n", err)
		return nil, ErrEmbeddingGenerationFailed
	}

	collection := rs.Mongodb.Database(knowledgeBase.Db).Collection(knowledgeBase.Collection)
	filter := bson.A{
		bson.D{
			{Key: "$vectorSearch",
				Value: bson.D{
					{Key: "index", Value: knowledgeBase.Index},
					{Key: "path", Value: knowledgeBase.Path},
					{Key: "queryVector", Value: queryEmbedding},
					{Key: "numCandidates", Value: knowledgeBase.Candidates},
					{Key: "limit", Value: knowledgeBase.Limit},
				},
			},
		},
	}
	cursor, err := collection.Aggregate(ctx, filter)
	if err != nil {
		log.Printf("Failed to find similar verses: %v\n", err)
		return nil, ErrSimilaritySearchFailed
	}
	defer cursor.Close(ctx)
	results := make([]interface{}, 0)
	for cursor.Next(ctx) {
		var data interface{}
		if err := cursor.Decode(&data); err != nil {
			return nil, fmt.Errorf("failed to decode data: %v", err)
		}
		results = append(results, data)
	}
	return results, nil
}

func (rs *RetrievalService) CreateChatCompletion(ctx context.Context, query Query) (string, error) {
	knowledgeBase, ok := rs.KnowledgeBases[query.Persona]
	if !ok {
		log.Printf("invalid persona %v\n", query.Persona)
		return "", ErrInvalidQueryParameters
	}

	documents, err := rs.SemanticSearch(ctx, query)
	if err != nil {
		return "", err
	}
	persona := knowledgeBase.Persona()
	contextString, err := persona.BuildCompletionContext(ctx, documents)
	if err != nil {
		log.Printf("Failed to build context: %v\n", err)
		return "", ErrChatCompletionContextBuildFailed
	}

	req := openai.ChatCompletionRequest{
		Model: openai.GPT4o,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: persona.Prompt(),
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: fmt.Sprintf("%s \n question: %s", contextString, query.Query),
			},
		},
	}

	response, err := rs.OpenAI.CreateChatCompletion(ctx, req)
	if err != nil {
		log.Printf("Error generating completion: %v", err)
		return "", ErrChatCompletionFailed
	}
	return response.Choices[0].Message.Content, nil
}
