package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq" // Регистрируем драйвер pq
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Enter search artist name")
	}
	// Строка подключения
	connStr := "user=SOFTuser dbname=soft password=SOFTpassword host=localhost sslmode=disable"

	// Открываем соединение. Имя драйвера - "postgres"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Проверяем соединение
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Успешное подключение")

	name := os.Args[1]
	albums, err := albumsByArtist(db, name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Albums found: %v\n", albums)

	id, err := strconv.ParseInt(os.Args[2], 10, 64)
	if err != nil {
		log.Fatal("ID must be a number:", err)
	}
	album, err := albumByID(db, id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Album found: %v\n", album)

	albID, err := addAlbum(db, Album{
		Title:  "The Modern Sound of Betty Carter",
		Artist: "Betty Carter",
		Price:  49.99,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID of added album: %v\n", albID)
}

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

func albumsByArtist(db *sql.DB, name string) ([]Album, error) {
	// An albums slice to hold data from returned rows.
	var albums []Album

	rows, err := db.Query("SELECT * FROM go.album WHERE artist = $1", name)
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	return albums, nil
}

func albumByID(db *sql.DB, id int64) (Album, error) {
	// An album to hold data from the returned row.
	var alb Album

	row := db.QueryRow("SELECT * FROM go.album WHERE id = $1", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("albumsById %d: no such album", id)
		}
		return alb, fmt.Errorf("albumsById %d: %v", id, err)
	}
	return alb, nil
}

func addAlbum(db *sql.DB, alb Album) (int64, error) {
	var id int64

	query := `INSERT INTO go.album (title, artist, price) 
              VALUES ($1, $2, $3) 
              RETURNING id`

	err := db.QueryRow(query, alb.Title, alb.Artist, alb.Price).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}

	return id, nil
}
