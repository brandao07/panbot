package todolist

import (
	"fmt"
	"time"
)

type category string

const (
	Movie      category = "Movie"
	TvShow     category = "TV Show"
	Anime      category = "Anime"
	Book       category = "Book"
	Song       category = "Song"
	MusicAlbum category = "Album"
)

var supportedCategories = map[category]bool{
	Movie:      true,
	TvShow:     true,
	Anime:      true,
	Book:       true,
	Song:       true,
	MusicAlbum: true,
}

type Item struct {
	Name        string    `json:"name"`
	Category    category  `json:"category"`
	AddedBy     string    `json:"added_by"`
	IsCompleted bool      `json:"is_completed"`
	CompletedAt time.Time `json:"completed_at"`
}

func NewItem(name string, category category, addedBy string) (*Item, error) {
	if !supportedCategories[category] {
		return nil, fmt.Errorf("category %s not supported", category)
	}

	item := &Item{
		Name:        name,
		Category:    category,
		AddedBy:     addedBy,
		IsCompleted: false,
		CompletedAt: time.Time{},
	}
	return item, nil
}
