package todolist

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

type Storage struct {
	filePath string
}

func NewStorage(filePath string) *Storage {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Println("File does not exist, creating new file")
		_, err := os.Create(filePath)
		if err != nil {
			panic(err)
		}
	}
	return &Storage{
		filePath: filePath,
	}
}

func (s *Storage) load() ([]Item, error) {
	file, err := os.Open(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("storage file does not exist")
		}
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("error closing storage file")
			panic(err)
		}
	}(file)

	var items []Item

	err = json.NewDecoder(file).Decode(&items)
	if err != nil {
		if err != io.EOF {
			return nil, err
		}
	}
	return items, nil
}

func (s *Storage) save(items []Item) error {
	file, err := os.Create(s.filePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("error closing storage file")
			panic(err)
		}
	}(file)

	err = json.NewEncoder(file).Encode(items)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) Add(item *Item) error {
	items, err := s.load()
	if err != nil {
		return err
	}

	// Check for duplicate names
	for _, v := range items {
		if v.Name == item.Name {
			return fmt.Errorf("item already exists")
		}
	}
	items = append(items, *item)
	return s.save(items)
}

func (s *Storage) Remove(item *Item) error {
	items, err := s.load()
	if err != nil {
		return err
	}
	for i, v := range items {
		if v.Name == item.Name {
			items = append(items[:i], items[i+1:]...)
			return s.save(items)
		}
	}
	return fmt.Errorf("item not found")
}

func (s *Storage) FindByName(name string) (*Item, error) {
	items, err := s.load()
	if err != nil {
		return nil, err
	}
	for _, v := range items {
		if v.Name == name {
			return &v, nil
		}
	}
	return nil, fmt.Errorf("item not found")
}

func (s *Storage) FindByCategory(category category) (*[]Item, error) {
	items, err := s.load()
	if err != nil {
		return nil, err
	}
	var filteredItems []Item
	for _, v := range items {
		if v.Category == category {
			filteredItems = append(filteredItems, v)
		}
	}
	return &filteredItems, nil
}

func (s *Storage) MarkAsCompleted(item *Item) error {
	items, err := s.load()
	if err != nil {
		return err
	}
	for i, v := range items {
		if v.Name == item.Name {
			items[i].IsCompleted = true
			items[i].CompletedAt = time.Now()
			return s.save(items)
		}
	}
	return fmt.Errorf("item not found")
}
