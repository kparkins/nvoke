package bible

import (
	"bufio"
	"fmt"
	"nvoke/pkg/embedding"
	"os"
	"strconv"
	"strings"
)

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

type VerseAdapter struct{}

func NewVerseAdapter() embedding.EmbeddingAdapter[*Verse] {
	return &VerseAdapter{}
}

func (vh *VerseAdapter) GetContent(verse *Verse) string {
	return verse.Text
}

func (vh *VerseAdapter) StoreEmbedding(verse *Verse, embedding []float32) {
	verse.Embedding = embedding
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

func ParseBook(line string) *Book {
	parts := strings.SplitN(line, " ", 2)
	return &Book{
		Name: strings.TrimSpace(parts[1]),
	}
}

func ParseVerse(line string) *Verse {
	parts := strings.SplitN(line, " ", 2)
	n, err := strconv.Atoi(parts[0])
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &Verse{
		Verse: n,
		Text:  strings.TrimSpace(parts[1]),
	}
}

func ParseBooksFile(f *os.File) []*Book {
	scanner := bufio.NewScanner(f)

	book := &Book{}
	verse := &Verse{}
	chapter := &Chapter{}
	books := make([]*Book, 0)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, " ", 2)

		// special case where the line has continued
		if !strings.HasPrefix(parts[0], "\\") {
			verse.Text = fmt.Sprintf("%s %s", verse.Text, strings.TrimSpace(line))
			continue
		}

		switch parts[0] {
		case "\\id":
			book = ParseBook(parts[1])
			books = append(books, book)
		case "\\c":
			n, err := strconv.Atoi(strings.TrimSpace(parts[1]))
			if err != nil {
				fmt.Println(err)
			}
			chapter = &Chapter{
				Number: n,
			}
			book.Chapters = append(book.Chapters, chapter)
		case "\\v":
			verse = ParseVerse(parts[1])
			verse.Book = book.Name
			verse.Chapter = chapter.Number
			chapter.Verses = append(chapter.Verses, verse)
		}
	}
	return books
}
