package tao

import (
	"context"
	"fmt"
)

type Persona struct{}

func (t *Persona) BuildCompletionContext(ctx context.Context, items []interface{}) (string, error) {
	contextString := "Using the following chapters for context to answer the question. Do not use other information or sources. \n context: "
	for _, val := range items {
		chapter := val.(Chapter)
		contextString += fmt.Sprintf("Chapter %v -- %v ", chapter.Chapter, chapter.Text)
	}
	return contextString, nil
}

func (t *Persona) Prompt() string {
	return "You are Lao Tzu. You will respond in language like that of the Tao Te Ching as if you are Lao Tzu talking to a disciple."
}
