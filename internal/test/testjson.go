package main

// import (
// 	"fmt"

// 	"github.com/jdetok/go-api-jdeko.me/internal/mariadb"
// 	"github.com/jdetok/go-api-jdeko.me/internal/store"
// )

// func main() {
// 	// store.RecGamesTest()

// 	db := mariadb.InitDB()
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
