package main

import (
	_  "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	"net/http"
	"log"
	"encoding/json"
	"io/ioutil"
	"bytes"
)

type Message struct {
	name  string
}

var db = &sql.DB{}

func getTodoList(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		w.Header().Set("Access-Control-Allow-Origin", "*");
		result,_ := query()
		fmt.Println(result)
		w.Write([]byte(result));

	} else if r.Method == "POST"{
		w.Header().Set("Access-Control-Allow-Origin", "*");
		req, _ := ioutil.ReadAll(r.Body)
		var data = bytes.NewBuffer(req).String()
		insert(data)
		fmt.Println()
	} else{
		fmt.Println(w);
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func insert(data string) {
	fmt.Println(data)
	stmt, err := db.Prepare(`INSERT todoList (name) values (?)`)
	checkErr(err)
	stmt.Exec(data)
}

func connectDatabase(){
	db,_= sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/todos?charset=utf8")
}

func query()(string,error){
	rows,_ :=db.Query("SELECT name FROM todoList")

	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}

	record := make([]map[string]interface{}, 0)

	for rows.Next(){
		rows.Scan(scanArgs...)
		entry := make(map[string]interface{})

		for i:=range values{
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry["name"] = v
		}

		record=append(record,entry)
	}
	result,_:=json.Marshal(record)

	return string(result),nil
}


func main() {
	connectDatabase()
	http.HandleFunc("/", getTodoList) //设置访问的路由
	err := http.ListenAndServe(":9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}