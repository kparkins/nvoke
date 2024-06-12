package bible

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Parser struct{}

func (p *Parser) Parse(location string) []*Verse {
	books := loadBooksFromFiles()
	return generateVerseDocuments(books)
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

func parseBooksFile(f *os.File) []*Book {
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

func loadBooksFromFiles() []*Book {
	books := make([]*Book, 0)
	directory := "texts/bible/nkjv/"
	entries, err := os.ReadDir(directory)
	if err != nil {
		fmt.Println(err)
		return books
	}
	for _, entry := range entries {
		f, err := os.Open(directory + entry.Name())
		if err != nil {
			fmt.Printf("Error opening file %v\n", err)
			continue
		}
		defer f.Close()
		books = append(books, parseBooksFile(f)...)
	}
	return books
}

func generateVerseDocuments(books []*Book) []*Verse {
	content := make([]*Verse, 0)
	for _, book := range books {
		for _, chapter := range book.Chapters {
			content = append(content, chapter.Verses...)
		}
	}
	return content
}
