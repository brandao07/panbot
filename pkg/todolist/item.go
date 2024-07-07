package todolist

import (
	"fmt"
	"strings"
	"time"
)

type category string

const (
	Movie      category = "Movie"
	TvShow     category = "TV Show"
	Anime      category = "Anime"
	Book       category = "Book"
	Song       category = "Song"
	MusicAlbum category = "Music Album"
)

var stringToCategory = map[string]category{
	"MOVIE":       Movie,
	"TV SHOW":     TvShow,
	"ANIME":       Anime,
	"BOOK":        Book,
	"SONG":        Song,
	"MUSIC ALBUM": MusicAlbum,
}

func convertStringToCategory(s string) (category, error) {
	s = strings.ToUpper(s)
	if cat, exists := stringToCategory[s]; exists {
		return cat, nil
	}
	return "", fmt.Errorf("invalid category")
}

type Item struct {
	Name        string    `json:"name"`
	Category    category  `json:"category"`
	AddedBy     string    `json:"added_by"`
	IsCompleted bool      `json:"is_completed"`
	CompletedAt time.Time `json:"completed_at"`
}

func NewItem(name string, category string, addedBy string) (*Item, error) {
	c, err := convertStringToCategory(category)
	if err != nil {
		return nil, err
	}

	item := &Item{
		Name:        name,
		Category:    c,
		AddedBy:     addedBy,
		IsCompleted: false,
		CompletedAt: time.Time{},
	}
	return item, nil
}

type ItemDTO struct {
	Name string
	User string
}

func NewItemDTO(name, user string) *ItemDTO {
	return &ItemDTO{
		Name: name,
		User: user,
	}
}
