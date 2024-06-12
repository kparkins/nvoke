package tao

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Parser struct{}

func (p *Parser) Parse(location string) []*Chapter {
	return LoadChaptersFromFile()
}

func ParseChapters(f *os.File) []*Chapter {
	scanner := bufio.NewScanner(f)

	var err error
	chapter := &Chapter{}
	chapters := make([]*Chapter, 0)

	for scanner.Scan() {
		line := scanner.Text()
		// special case where the line has continued
		if strings.HasPrefix(line, "Chapter") {
			chapters = append(chapters, chapter)
			chapter = &Chapter{}
			parts := strings.SplitN(line, " ", 2)
			chapter.Chapter, err = strconv.Atoi(parts[1])
			if err != nil {
				log.Printf("Error parsing chapter number: %v", err)
				return nil
			}
		} else {
			chapter.Text = fmt.Sprintf("%s %s", chapter.Text, strings.TrimSpace(line))
		}

	}
	return chapters
}

func LoadChaptersFromFile() []*Chapter {
	f, err := os.Open("texts/tao/linn/tao.txt")
	if err != nil {
		fmt.Printf("Failed to read the tao source text: %v\n", err)
		return nil
	}
	chapters := ParseChapters(f)
	return chapters
}
