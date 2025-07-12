package src

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Root struct {
	PacketVersion               int             `json:"packetVersion"`
	FromVersion                 int             `json:"fromVersion"`
	CatalogTablesVersion        int             `json:"catalogTablesVersion"`
	CatalogSpecialTablesVersion int             `json:"catalogSpecialTablesVersion"`
	CatalogEventViewVersion     int             `json:"catalogEventViewVersion"`
	SportBasicMarketsVersion    int             `json:"sportBasicMarketsVersion"`
	SportBasicFactorsVersion    int             `json:"sportBasicFactorsVersion"`
	IndependentFactorsVersion   int             `json:"independentFactorsVersion"`
	FactorsVersion              int             `json:"factorsVersion"`
	ComboFactorsVersion         int             `json:"comboFactorsVersion"`
	SportKindsVersion           int             `json:"sportKindsVersion"`
	TopCompetitionsVersion      int             `json:"topCompetitionsVersion"`
	EventSmartFiltersVersion    int             `json:"eventSmartFiltersVersion"`
	GeoCategoriesVersion        int             `json:"geoCategoriesVersion"`
	SportCategoriesVersion      int             `json:"sportCategoriesVersion"`
	TournamentInfos             []any           `json:"tournamentInfos"`
	Sports                      []Sport         `json:"sports"`
	Events                      []Event         `json:"events"`
	EventBlocks                 []any           `json:"eventBlocks"`
	EventMiscs                  []EventMisc     `json:"eventMiscs"`
	LiveEventInfos              []LiveEventInfo `json:"liveEventInfos"`
	CustomFactors               []CustomFactor  `json:"customFactors"`
}

type Sport struct {
	ID               int    `json:"id"`
	Kind             string `json:"kind"`
	SortOrder        string `json:"sortOrder"`
	Name             string `json:"name"`
	Alias            string `json:"alias,omitempty"`
	ParentID         int    `json:"parentId,omitempty"`
	ParentIDs        []int  `json:"parentIds,omitempty"`
	RegionID         int    `json:"regionId,omitempty"`
	TournamentInfoID int    `json:"tournamentInfoId,omitempty"`
	SportCategoryID  int    `json:"sportCategoryId,omitempty"`
}

type Event struct {
	ID        int    `json:"id"`
	SortOrder string `json:"sortOrder"`
	Level     int    `json:"level"`
	Num       int    `json:"num"`
	SportID   int    `json:"sportId"`
	Kind      int    `json:"kind"`
	RootKind  int    `json:"rootKind"`
	Team1ID   int    `json:"team1Id"`
	Team2ID   int    `json:"team2Id"`
	Team1     string `json:"team1"`
	Team2     string `json:"team2"`
	Name      string `json:"name"`
	StartTime int64  `json:"startTime"`
	Place     string `json:"place"`
	Priority  int    `json:"priority"`
	TV        []int  `json:"tv"`
}

type EventMisc struct {
	ID             int    `json:"id"`
	LiveDelay      int    `json:"liveDelay"`
	Score1         int    `json:"score1"`
	Score2         int    `json:"score2"`
	Comment        string `json:"comment"`
	TimerDirection int    `json:"timerDirection"`
	TimerSeconds   int    `json:"timerSeconds"`
}

type LiveEventInfo struct {
	EventID        int       `json:"eventId"`
	Timer          string    `json:"timer"`
	TimerSeconds   int       `json:"timerSeconds"`
	TimerDirection int       `json:"timerDirection"`
	ScoreFunction  string    `json:"scoreFunction"`
	ScoreComment   string    `json:"scoreComment"`
	Scores         [][]Score `json:"scores"`
	Subscores      []any     `json:"subscores"`
}

type Score struct {
	C1    string `json:"c1"`
	C2    string `json:"c2"`
	Title string `json:"title,omitempty"`
}

type CustomFactor struct {
	E        int      `json:"e"`
	CountAll int      `json:"countAll"`
	Factors  []Factor `json:"factors"`
}

type Factor struct {
	F  int     `json:"f"`
	V  float64 `json:"v"`
	P  *int    `json:"p,omitempty"`
	PT *string `json:"pt,omitempty"`
}

type Request struct {
	Game_id string `json:"game_id"`
}

func Get_json_file(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request Request
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			log.Fatal("Ошибка в декодировании http запроса:", err)
		}

		req_Url := fmt.Sprintf("https://line31w.bk6bba-resources.com/events/event?lang=en&version=55248127755&eventId=%s&scopeMarket=1600", request.Game_id)

		resp, err := http.Get(req_Url)
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

		// Вставляем матч в таблицу matches
		if _, err := db.Exec(
			"INSERT INTO matches(game_id, url, count) VALUES ($1, $2, $3)",
			request.Game_id, req_Url, jsons.EventMiscs[0].Comment); err != nil {
			log.Fatal("Ошибка при вставке в matches:", err)
		}

		tableName := request.Game_id
		createTableQuery := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "%s" (
			number INT,
			value  DOUBLE PRECISION
		)`, tableName)

		if _, err := db.Exec(createTableQuery); err != nil {
			log.Fatal("Ошибка при создании таблицы:", tableName, err)
		}

		insertQuery := fmt.Sprintf(`INSERT INTO "%s" (number, value) VALUES ($1, $2)`, tableName)
		for _, cf := range jsons.CustomFactors {
			for _, f := range cf.Factors {
				_, err := db.Exec(insertQuery, f.F, f.V)
				if err != nil {
					log.Println("Ошибка при вставке фактора:", err)
				}
			}
		}
	}
}
