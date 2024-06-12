package bible

import "nvoke/pkg/embedding"

type EmbeddingAdapter struct{}

func NewEmbeddingAdapter() embedding.Adapter[*Verse] {
	return &EmbeddingAdapter{}
}

func (vh *EmbeddingAdapter) GetContent(verse *Verse) string {
	return verse.Text
}

func (vh *EmbeddingAdapter) StoreEmbedding(verse *Verse, embedding []float32) {
	verse.Embedding = embedding
}
