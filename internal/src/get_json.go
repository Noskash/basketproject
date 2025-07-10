package src

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func Get_json_file(game_id string) string {
	reqUrl := fmt.Sprintf("https://line31w.bk6bba-resources.com/events/event?lang=en&version=55248127755&eventId=%s&scopeMarket=1600", game_id)

	resp, err := http.Get(reqUrl)

	if err != nil {
		log.Fatal("Ошибка в отпрапвке GET запроса для получения json", err)
	}

	respBody, err := io.ReadAll(resp.Body)

	resp.Body.Close()

	return string(respBody)
}
