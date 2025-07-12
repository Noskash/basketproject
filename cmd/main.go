package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Noskash/basketproject/db"
	"github.com/Noskash/basketproject/internal/src"
)

func main() {
	db, err := db.Connect_to_database()

	if err != nil {
		log.Fatal("Ошибка при подключении к бд")
	}

	http.HandleFunc("/addmatch", src.Get_json_file(db))

	ticker := time.NewTicker(5 * time.Second)

	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			src.Update_values(db)
		}
	}
}
