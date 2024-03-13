package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// Item represents a simple item with ID and Name.
type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func createTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL
	);`
	_, err := db.Exec(query)
	return err
}

func insertItem(db *sql.DB, item Item) (int64, error) {
	result, err := db.Exec("INSERT INTO items (name) VALUES (?)", item.Name)
	if err != nil {
		return 0, err
	}
	id, _ := result.LastInsertId()
	return id, nil
}

func getItem(db *sql.DB, id int) (Item, error) {
	var item Item
	err := db.QueryRow("SELECT id, name FROM items WHERE id=?", id).Scan(&item.ID, &item.Name)
	return item, err
}

func updateItem(db *sql.DB, item Item) error {
	_, err := db.Exec("UPDATE items SET name=? WHERE id=?", item.Name, item.ID)
	return err
}

func deleteItem(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM items WHERE id=?", id)
	return err
}

func getAllItems(db *sql.DB) ([]Item, error) {
	rows, err := db.Query("SELECT id, name FROM items")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]Item, 0)
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ID, &item.Name)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func main() {
	db, err := sql.Open("sqlite3", "./items.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = createTable(db)
	if err != nil {
		log.Fatal(err)
	}

	item := Item{Name: "Sample Item"}
	id, err := insertItem(db, item)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted item with ID: %d\n", id)

	itemFromDB, err := getItem(db, int(id))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Retrieved item from DB: %+v\n", itemFromDB)

	itemFromDB.Name = "Updated Item"
	err = updateItem(db, itemFromDB)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Item updated successfully")

	allItems, err := getAllItems(db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("All items in the DB: %+v\n", allItems)

	err = deleteItem(db, int(id))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Item deleted successfully")
}
