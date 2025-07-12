package src

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func Update_values(db *sql.DB) {
	urls, err := db.Query("SELECT game_id ,url FROM matches")

	if err != nil {
		log.Fatal("Ошибка во время выборки url из бд", err)
	}

	defer urls.Close()

	for urls.Next() {
		var url string
		var game_id string
		err := urls.Scan(&game_id, &url)

		if err != nil {
			log.Fatal("Ошибка при сканировании данных из бд")
		}

		resp, err := http.Get(url)

		if err != nil {
			log.Fatal("Ошибка в отпрапвке GET запроса для получения json", err)
		}

		respBody, err := io.ReadAll(resp.Body)

		if err != nil {
			log.Fatal("Ошибка при декодировании ответа от fonbet")
		}

		resp.Body.Close()

		var jsons Root

		if err := json.Unmarshal(respBody, &jsons); err != nil {
			log.Fatal("Ошибка при декодировании json")
		}
		new_title := game_id + "_new"
		for _, cf := range jsons.CustomFactors {
			for _, fv := range cf.Factors {
				var exists bool

				if err := db.QueryRow("SELECT EXISTS(*) FROM %s WHERE number = $1", game_id, fv.F).Scan(&exists); err != nil {
					log.Fatal("Ошибка при запросе в базу данных ", game_id)
				}
				if !exists {
					_, err := db.Exec("INSERT INTO %s (number , value) VALUES($1 , $2)", game_id, fv.F, fv.V)
					if err != nil {
						log.Fatal("Ошибка при вставке данных в ", game_id)
					}
				} else {
					_, err := db.Exec("DROP TABLE IF EXISTS %s", new_title)

					if err != nil {
						log.Fatal("Ошибка при удалении временной бд номером", game_id)
					}

					if _, err := db.Exec("CREATE TABLE %s(number int , value int)", new_title); err != nil {
						log.Fatal("Ошибка при создании временной бд с номером ", new_title, err)
					}

					if _, err := db.Exec("INSERT INTO %s(number , value) VALUES ($1 , $2)", new_title, fv.F, fv.V); err != nil {
						log.Fatal("Ошибка при вставке данных во временную базу данных с номером", new_title, err)
					}
				}
			}
		}
		if _, err := db.Exec("DROP TABLE IF EXISTS %s", game_id); err != nil {
			log.Fatal("Ошибка при удалении таблицы с номером", game_id)
		}

		if _, err := db.Exec("ALTER TABLE %s RENAME $1", new_title, game_id); err != nil {
			log.Fatal("Ошибка при замене базы данных с номером", game_id)
		}
	}
}
