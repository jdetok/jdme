package main

// import (
// 	"encoding/json"
// 	"fmt"

// 	"github.com/jdetok/go-api-jdeko.me/internal/mariadb"
// 	"github.com/jdetok/go-api-jdeko.me/internal/store"
// )

// func main() {
// 	// store.RecGamesTest()

// 	db := mariadb.InitDB()
// 	var rp store.Resp

// 	rp.GetPlayerDash(db, 2544, 42024)
// 	js, err := json.Marshal(rp)
// 	if err != nil {
// 		fmt.Println("marshal error")
// 	}

// 	fmt.Println(string(js))

// rows, err := db.Query(`
// 	select * from api_player_stats
// 	where player_id = ? and season_id = ?
// 	`, 2544, 22005)
// if err != nil {
// 	fmt.Println("query error")
// }

// cols, _ := rows.Columns()
// _, err = mariadb.ProcessRows(rows, cols)
// if err != nil {
// 	fmt.Println("row processing error")
// }

// }

// 	rg := store.RecentGames{}
// 	ts := store.TopScorer{}

// 	js, err := rg.Get(db)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println(string(js))

// 	jsTs, err := ts.Get(db)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println(string(jsTs))

// 	// fmt.Println(rg.GameId)
// 	// fmt.Println(rg.GameDate)
// 	// fmt.Println(rg.Final)
// 	// fmt.Println(rg.Overtime)

// 	// rows, err := db.Query(mariadb.RecentGames.Q)
// 	// if err != nil {
// 	// 	fmt.Println(err)
// 	// }

// 	// rg.Make(rows)

// 	// js, err := json.Marshal(rg)
// 	// if err != nil {
// 	// 	fmt.Println(err)
// 	// }

// 	// fmt.Println(string(js))

// 	// results, err := mariadb.ProcessRows(rows, cols)
// 	// if err != nil {
// 	// 	fmt.Println(err)
// 	// }

// 	// for i, r := range results {

// 	// }

// // old testing:
// 	// res := jsonops.MapJSONFile("json/teams.json")
// 	// // fmt.Println(res)

// 	// var body []byte = jsonops.MapToJSON("", res)

// 	// jsonops.SaveJSON("json/test.json", body)
// }
