package src

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func Update_values(db *sql.DB) {
	urls, err := db.Query("SELECT game_id, url FROM matches")
	if err != nil {
		log.Fatal("Ошибка во время выборки url из БД:", err)
	}
	defer urls.Close()

	for urls.Next() {
		var gameID string
		var url string
		if err := urls.Scan(&gameID, &url); err != nil {
			log.Fatal("Ошибка при сканировании данных из БД:", err)
		}

		resp, err := http.Get(url)
		if err != nil {
			log.Fatal("Ошибка при GET-запросе:", err)
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("Ошибка при чтении тела ответа:", err)
		}

		var jsons Root
		if err := json.Unmarshal(respBody, &jsons); err != nil {
			log.Fatal("Ошибка при декодировании JSON:", err)
		}

		newTable := gameID + "_new"

		createNewTable := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "%s" (
			number INT,
			value  DOUBLE PRECISION
		)`, newTable)
		if _, err := db.Exec(createNewTable); err != nil {
			log.Fatal("Ошибка при создании временной таблицы:", err)
		}

		insertNew := fmt.Sprintf(`INSERT INTO "%s" (number, value) VALUES ($1, $2)`, newTable)
		for _, cf := range jsons.CustomFactors {
			for _, fv := range cf.Factors {
				_, err := db.Exec(insertNew, fv.F, fv.V)
				if err != nil {
					log.Println("Ошибка при вставке значения в", newTable, ":", err)
				}
			}
		}

		dropOld := fmt.Sprintf(`DROP TABLE IF EXISTS "%s"`, gameID)
		if _, err := db.Exec(dropOld); err != nil {
			log.Fatal("Ошибка при удалении старой таблицы:", err)
		}

		rename := fmt.Sprintf(`ALTER TABLE "%s" RENAME TO "%s"`, newTable, gameID)
		if _, err := db.Exec(rename); err != nil {
			log.Fatal("Ошибка при переименовании временной таблицы:", err)
		}
	}
}
