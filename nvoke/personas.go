package nvoke

import (
	"context"
	"nvoke/pkg/bible"
	"nvoke/pkg/tao"
)

type Persona interface {
	BuildCompletionContext(ctx context.Context, items []interface{}) (string, error)
	Prompt() string
}

type KnowledgeBase struct {
	Index      string
	Path       string
	Db         string
	Collection string
	Limit      int
	Candidates int
	persona    Persona
}

func (kb *KnowledgeBase) Persona() Persona {
	return kb.persona
}

var KnowledgeBases = map[string]KnowledgeBase{
	"bible": {
		Index:      "embedding",
		Path:       "embedding",
		Db:         "bible",
		Collection: "verses",
		Limit:      20,
		Candidates: 200,
		persona:    &bible.Persona{},
	},
	"tao": {
		Index:      "embedding",
		Path:       "embedding",
		Db:         "tao",
		Collection: "chapters",
		Limit:      20,
		Candidates: 200,
		persona:    &tao.Persona{},
	},
}
