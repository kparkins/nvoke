package tao

type Chapter struct {
	Chapter   int       `json:"chapter"`
	Text      string    `json:"text"`
	Embedding []float32 `json:"embedding,omitempty"`
}
