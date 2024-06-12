package tao

import "nvoke/pkg/embedding"

type EmbeddingAdapter struct{}

func NewEmbeddingAdapter() embedding.Adapter[*Chapter] {
	return &EmbeddingAdapter{}
}

func (a *EmbeddingAdapter) GetContent(chapter *Chapter) string {
	return chapter.Text
}

func (a *EmbeddingAdapter) StoreEmbedding(chapter *Chapter, embedding []float32) {
	chapter.Embedding = embedding
}
