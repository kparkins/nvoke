package bible

type Book struct {
	Name     string     `json:"book_name"`
	Chapters []*Chapter `json:"chapters"`
}

type Chapter struct {
	Number int      `json:"chapter_number"`
	Verses []*Verse `json:"verses"`
}

type Verse struct {
	Book      string    `json:"book"`
	Chapter   int       `json:"chapter"`
	Verse     int       `json:"verse"`
	Text      string    `json:"text"`
	Embedding []float32 `json:"embedding,omitempty"`
}

func NewBook(name string) *Book {
	return &Book{
		Name: name,
	}
}

func NewChapter(number int) *Chapter {
	return &Chapter{
		Number: number,
	}
}

func NewVerse(chapter int, verse int, text string) *Verse {
	return &Verse{
		Chapter: chapter,
		Verse:   verse,
		Text:    text,
	}
}
