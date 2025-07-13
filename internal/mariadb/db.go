package mariadb

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jdetok/go-api-jdeko.me/internal/env"
	"github.com/jdetok/go-api-jdeko.me/internal/errs"
)

type JSONOutput struct {
	Meta []string         `json:"meta"`
	Data []map[string]any `json:"data"`
}

func InitDB() *sql.DB {
	// get conn. vars from .env & build connection string
	dbUser := env.GetString("DB_USER")
	dbHost := env.GetString("DB_HOST")
	database := env.GetString("DB")
	connStr := dbUser + "@tcp(" + dbHost + ")/" + database

	// attempt to open database connection
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		fmt.Println(err)
	}

	// ping to confirm connection - exits with error if ping is not successful
	if err := db.Ping(); err != nil {
		fmt.Println(err)
	}
	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(200)
	return db
}

// run funcs below to query db and return a json []byte
func DBJSONResposne(db *sql.DB, q string, args ...any) ([]byte, error) {
	rows, cols, err := Select(db, q, args...)
	if err != nil {
		fmt.Printf("Error occured querying db: %v\n", err)
		return nil, err
	}
	// returns slice of any
	vals, err := ProcessRows(rows, cols)
	if err != nil {
		fmt.Printf("Error occured processing rows: %v\n", err)
		return nil, err
	}
	// returns map to marshal
	mapData, err := MapRows(vals, cols)
	if err != nil {
		fmt.Printf("Error occured converting rows to map: %v\n", err)
		return nil, err
	}
	// returns json []byte
	j, err := MapToJSON(mapData)
	if err != nil {
		fmt.Printf("Error occured marshalling map: %v\n", err)
		return nil, err
	}
	return j, err
}

// new database query func, use variatic args
func Select(db *sql.DB, q string, args ...any) (*sql.Rows, []string, error) {
	// actual query
	rows, err := db.Query(q, args...)
	if err != nil {
		fmt.Println("failed to query db")
		return nil, nil, err
	}

	// slice of column names
	cols, err := rows.Columns()
	if err != nil {
		fmt.Println("failed to get columns")
		return nil, nil, err
	}
	return rows, cols, nil
}

func PrintRows(rows *sql.Rows) {
	for rows.Next() {

	}
}

// process rows returned by Select function
func ProcessRows(rows *sql.Rows, cols []string) ([]any, error) {
	var vals []any // final slice with vals
	for rows.Next() {
		// make slice of non-nil pointers to scan
		ptrs := make([]any, len(cols))
		for i := range ptrs {
			ptrs[i] = new(any)
		}

		// scan for the ptrs
		if err := rows.Scan(ptrs...); err != nil {
			fmt.Println("couldn't scan")
		}

		// make an any slice to deref ptrs for each row
		rowVals := make([]any, len(ptrs))
		for i := range ptrs {
			v := *(ptrs[i].(*any)) // extract val from each ptr

			// convert byte slices to strings
			if reflect.TypeOf(v) == reflect.TypeOf([]byte{}) {
				v = string(v.([]byte))
			}

			rowVals[i] = v
		}
		// append row vals to vals slice
		vals = append(vals, rowVals)

	}
	fmt.Println(vals)
	return vals, nil
}

// use cols and vals to make key-val pairs and marshal to json
func MapRows(vals []any, cols []string) ([]map[string]any, error) {
	var data []map[string]any
	for _, val := range vals {
		val, ok := val.([]any) // force []any type to be able to iterate
		if !ok {
			return nil, errors.New("error forcing []any type")
		}

		// create a map per val with column and val
		valMap := map[string]any{}
		for i, v := range val {
			valMap[cols[i]] = v // key value pair for each column
		}
		// append to slice of maps
		data = append(data, valMap)
	}
	return data, nil
}

// marshal slice of maps to JSON
func MapToJSON(data []map[string]any) ([]byte, error) {
	jsonB, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Erroring marshalling: %v", err)
	}
	return jsonB, nil
}

// =======================================================================

func SelectDB(database *sql.DB, q string) ([]byte, error) {
	e := errs.ErrInfo{Prefix: "database query"}
	rows, err := database.Query(q)
	if err != nil {
		fmt.Println(err)
	}
	js, err := RowsToJSON(rows, false)
	if err != nil {
		e.Msg = "func RowsToJSON() failed"
		return nil, e.Error(err)
	}
	// return the response as json
	return js, nil
}

// USED IN ORIGINAL NBA ENDPOINTS
func SelectLgPlayer(db *sql.DB, w *http.ResponseWriter, q string, lg string, pl string) []byte {
	e := errs.ErrInfo{Prefix: "database query (arg)"}

	rows, err := db.Query(q, lg, pl)
	if err != nil {
		e.Msg = "db.Query failed"
		errs.HTTPErr(*w, e.Error(err))
	}

	// return the response as json
	js, err := RowsToJSON(rows, false)
	if err != nil {
		e.Msg = "func RowsToJSON() failed"
		errs.HTTPErr(*w, e.Error(err))
	}
	return js
}
