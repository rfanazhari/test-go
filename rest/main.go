package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

// Item represents a simple item with ID and Name.
type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// ItemStore handles storage and retrieval of items.
type ItemStore struct {
	sync.Mutex
	items     map[int]Item
	currentID int
}

// NewItemStore creates a new ItemStore.
func NewItemStore() *ItemStore {
	return &ItemStore{
		items:     make(map[int]Item),
		currentID: 1,
	}
}

func (store *ItemStore) getItemsHandler(w http.ResponseWriter, r *http.Request) {
	store.Lock()
	defer store.Unlock()

	items := make([]Item, 0, len(store.items))
	for _, item := range store.items {
		items = append(items, item)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func (store *ItemStore) postItemHandler(w http.ResponseWriter, r *http.Request) {
	store.Lock()
	defer store.Unlock()

	var newItem Item
	err := json.NewDecoder(r.Body).Decode(&newItem)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error decoding request body: %v", err), http.StatusBadRequest)
		return
	}

	newItem.ID = store.currentID
	store.currentID++
	store.items[newItem.ID] = newItem

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newItem)
}

func main() {
	itemStore := NewItemStore()

	http.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			itemStore.getItemsHandler(w, r)
		case http.MethodPost:
			itemStore.postItemHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	port := 8080
	fmt.Printf("Starting server on :%d...\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
