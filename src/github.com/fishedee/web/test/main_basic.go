package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"strconv"
	"time"
)

var db *sql.DB

func init() {
	db, _ = sql.Open("mysql", "root:1@tcp(127.0.0.1:3306)/bakeweb?charset=utf8")
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)
	db.Ping()
}

func dbQuery(sql string) []map[string]string {
	rows, err := db.Query(sql)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		panic(err)
	}
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}

	result := []map[string]string{}
	for rows.Next() {
		//将行数据保存到record字典
		record := map[string]string{}
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err)
		}
		for i, col := range values {
			if col != nil {
				record[columns[i]] = string(col.([]byte))
			}
		}
		result = append(result, record)
	}
	return result
}

func jsonOutput(w http.ResponseWriter, data interface{}) {
	var result struct {
		Code int         `json:"code"`
		Msg  string      `json:"msg"`
		Data interface{} `json:"data"`
	}
	result.Data = data
	dataByte, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}
	w.Write(dataByte)
}

func Doing(w http.ResponseWriter, r *http.Request) {
	jsonOutput(w, "Hello World")
}

func DbTask(w http.ResponseWriter, r *http.Request) {
	var result struct {
		Count int                 `json:"count"`
		Data  []map[string]string `json:"data"`
	}
	countArray := dbQuery("select count(*) from t_user")
	count, err := strconv.Atoi(countArray[0]["count(*)"])
	if err != nil {
		panic(err)
	}

	data := dbQuery("select * from t_user")

	result.Count = count
	result.Data = data

	jsonOutput(w, result)
}

func LongTask(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("task begin: %v\n", time.Now())
	time.Sleep(time.Second * 5)
	fmt.Printf("task end: %v\n", time.Now())
	jsonOutput(w, nil)
}

func main() {
	http.HandleFunc("/doing", Doing)
	http.HandleFunc("/dbtask", DbTask)
	http.HandleFunc("/longtask", LongTask)
	fmt.Println("Server Run in 9011!")
	err := http.ListenAndServe(":9011", nil)
	if err != nil {
		panic(err)
	}
}
