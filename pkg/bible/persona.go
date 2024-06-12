package bible

import (
	"context"
	"fmt"
)

type Persona struct{}

func (b *Persona) BuildCompletionContext(ctx context.Context, items []interface{}) (string, error) {
	contextString := "Using the following verses for context to answer the question. Do not use other information or sources. \n context: "
	for _, val := range items {
		verse := val.(Verse)
		contextString += fmt.Sprintf("%v %v:%v -- %v ", verse.Book, verse.Chapter, verse.Verse, verse.Text)
	}
	return contextString, nil
}

func (b *Persona) Prompt() string {
	return "You are Jesus. You will respond in language like that of the NKJV bible as if you are Jesus talking to his son."
}
