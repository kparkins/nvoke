package nvoke

import "github.com/sashabaranov/go-openai"

type ChatCompletionStreamResponse openai.ChatCompletionStreamResponse

type StreamingResponse[T any] interface {
	Recv() (T, error)
	Close() error
}

type ChatCompletionStreamAdapter struct {
	stream *openai.ChatCompletionStream
}

func (a *ChatCompletionStreamAdapter) Recv() (ChatCompletionStreamResponse, error) {
	response, err := a.stream.Recv()
	if err != nil {
		return ChatCompletionStreamResponse{}, err
	}
	return ChatCompletionStreamResponse(response), nil
}

func (a *ChatCompletionStreamAdapter) Close() error {
	return a.stream.Close()
}
