package todolist

import (
	"os"
	"reflect"
	"testing"
)

func TestStorage_Add(t *testing.T) {
	// Create a temporary test file
	tmpFile := createTempFile(t)
	storage := NewStorage(tmpFile)

	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			panic(err)
		}
	}(tmpFile)

	// Test adding a new item
	item1, err := NewItem("Item1", Movie, "André")
	if err != nil {
		t.Errorf("Error creating item: %s", err)
	}
	err = storage.Add(item1)
	if err != nil {
		t.Errorf("Error adding item: %v", err)
	}

	// Test adding a duplicate item
	err = storage.Add(item1)
	if err == nil {
		t.Errorf("Expected error adding duplicate item, but got none")
	}

	// Test adding another item
	item2, err := NewItem("Item2", Song, "André")
	if err != nil {
		t.Errorf("Error creating item: %s", err)
	}

	err = storage.Add(item2)
	if err != nil {
		t.Errorf("Error adding item: %v", err)
	}

	// Test if items were added correctly
	items, err := storage.load()
	if err != nil {
		t.Errorf("Error loading items: %v", err)
	}

	expectedItems := []Item{*item1, *item2}
	if !reflect.DeepEqual(items, expectedItems) {
		t.Errorf("Loaded items do not match expected items. Got: %v, Expected: %v", items, expectedItems)
	}
}

func TestStorage_FindByName(t *testing.T) {
	// Create a temporary test file
	tmpFile := createTempFile(t)
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			panic(err)
		}
	}(tmpFile)

	storage := NewStorage(tmpFile)

	// Add an item
	item, err := NewItem("Item", Movie, "André")
	if err != nil {
		t.Errorf("Error creating item: %s", err)
	}
	_ = storage.Add(item)

	// Test finding an existing item
	foundItem, err := storage.FindByName("Item")
	if err != nil {
		t.Errorf("Error finding item: %v", err)
	}
	if foundItem == nil || foundItem.Name != "Item" {
		t.Errorf("Expected to find item 'Item', but found: %v", foundItem)
	}

	// Test finding a non-existing item
	_, err = storage.FindByName("Non-existing item")
	if err == nil {
		t.Errorf("Expected error finding non-existing item, but got none")
	}
}

func TestStorage_MarkAsCompleted(t *testing.T) {
	// Create a temporary test file
	tmpFile := createTempFile(t)
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			panic(err)
		}
	}(tmpFile)

	storage := NewStorage(tmpFile)

	// Add an item
	item, err := NewItem("Item", Movie, "André")
	if err != nil {
		t.Errorf("Error creating item: %s", err)
	}
	_ = storage.Add(item)

	// Mark the item as completed
	err = storage.MarkAsCompleted(item)
	if err != nil {
		t.Errorf("Error marking item as completed: %v", err)
	}

	// Check if the item was marked as completed correctly
	items, _ := storage.load()
	if len(items) != 1 || !items[0].IsCompleted {
		t.Errorf("Item was not marked as completed correctly")
	}
}

func TestStorage_Remove(t *testing.T) {
	tmpFile := createTempFile(t)
	storage := NewStorage(tmpFile)

	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			panic(err)
		}
	}(tmpFile)

	// Add an item
	item, err := NewItem("Item", Movie, "André")
	if err != nil {
		t.Errorf("Error creating item: %s", err)
	}
	_ = storage.Add(item)

	// Remove Item
	err = storage.Remove(item)
	if err != nil {
		t.Errorf("Error removing item: %v", err)
	}

	// Check if Item exists still
	_, err = storage.FindByName(item.Name)
	if err == nil {
		t.Errorf("Expected error finding non-existing item, but got none")
	}
}

func createTempFile(t *testing.T) string {
	tmpFile, err := os.Create("test.json")
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}
	err = tmpFile.Close()
	if err != nil {
		panic(err)
	}

	return tmpFile.Name()
}
