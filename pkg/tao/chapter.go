package tao

import "nvoke/pkg/embedding"

type Chapter struct {
	Chapter   int       `json:"chapter"`
	Text      string    `json:"text"`
	Embedding []float32 `json:"embedding,omitempty"`
}

type Adapter struct{}

func NewAdapter() embedding.EmbeddingAdapter[*Chapter] {
	return &Adapter{}
}

func (a *Adapter) GetContent(chapter *Chapter) string {
	return chapter.Text
}

func (a *Adapter) StoreEmbedding(chapter *Chapter, embedding []float32) {
	chapter.Embedding = embedding
}
