package main

import (
	"fmt"

	"github.com/jdetok/go-api-jdeko.me/internal/mariadb"
)

func main() {
	db := mariadb.InitDB()
	// q := `
	// 	select * from v_szn_avgs
	// 	where season = ?
	// 	and team = ?
	// `

	rows, cols, err := mariadb.Select(db, "select * from career_avgs")
	if err != nil {
		fmt.Printf("Error occured querying db: %v\n", err) 
	}
	
	resp, err := mariadb.ProcessRows(rows, cols)
	fmt.Println(resp)

	js, err := mariadb.RowsToJSON(rows, false)
	if err != nil {
		fmt.Printf("Error with rows to json function: %v\n", err) 
	}

	fmt.Println(string(js))
	// resp, err := mariadb.DBJSONResposne(db, "select * from career_avgs")
	// if err != nil {
	// 	fmt.Printf("Error occured querying db: %v\n", err)
	// }

	// fmt.Println(string(resp))

	// var data any
	// json.Unmarshal(resp, &data)
	// fmt.Println(data)
}
