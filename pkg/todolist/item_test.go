package todolist

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestNewItem(t *testing.T) {
	tests := []struct {
		name        string
		category    category
		addedBy     string
		expectedErr bool
	}{
		{
			name:        "Test Movie Category",
			category:    Movie,
			addedBy:     "UserA",
			expectedErr: false,
		},
		{
			name:        "Test Unsupported Category",
			category:    "UnsupportedCategory",
			addedBy:     "UserB",
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item, err := NewItem("Example Item", tt.category, tt.addedBy)

			if tt.expectedErr {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if item == nil {
				t.Fatal("Item is nil but no error returned")
			}

			if item.Name != "Example Item" {
				t.Errorf("Expected item name to be 'Example Item', got '%s'", item.Name)
			}

			if item.Category != tt.category {
				t.Errorf("Expected category to be '%s', got '%s'", tt.category, item.Category)
			}

			if item.AddedBy != tt.addedBy {
				t.Errorf("Expected added by to be '%s', got '%s'", tt.addedBy, item.AddedBy)
			}

			if item.IsCompleted {
				t.Error("Expected item to be not completed, but it is marked as completed")
			}

			if !item.CompletedAt.IsZero() {
				t.Error("Expected completedAt to be zero value, but it is not")
			}
		})
	}
}

func TestNewItemJSONMarshalling(t *testing.T) {
	item := &Item{
		Name:        "Example Item",
		Category:    Movie,
		AddedBy:     "UserA",
		IsCompleted: false,
		CompletedAt: time.Time{},
	}

	jsonData := `{"name":"Example Item","category":"Movie","added_by":"UserA","is_completed":false,"completed_at":"0001-01-01T00:00:00Z"}`

	t.Run("Marshal Item to JSON", func(t *testing.T) {
		b, err := json.Marshal(item)
		if err != nil {
			t.Fatalf("Error marshalling item to JSON: %v", err)
		}
		if string(b) != jsonData {
			t.Errorf("Expected marshalled JSON to be '%s', got '%s'", jsonData, string(b))
		}
	})

	t.Run("Unmarshal JSON to Item", func(t *testing.T) {
		var newItem Item
		err := json.Unmarshal([]byte(jsonData), &newItem)
		if err != nil {
			t.Fatalf("Error unmarshalling JSON to item: %v", err)
		}
		if !reflect.DeepEqual(newItem, *item) {
			t.Errorf("Expected unmarshalled item to be equal to original item, but they are not")
		}
	})
}
