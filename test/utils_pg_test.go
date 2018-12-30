package main

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/lib/pq"
)

func Test_pg(t *testing.T) {
	db, err := sql.Open("postgres", "user=ttank password=ttank123 dbname=ttank sslmode=disable")
	checkErr(err)

	//插入数据
	stmt, err := db.Prepare("INSERT INTO person(id,name,status) VALUES($1,$2,$3) RETURNING id")
	checkErr(err)

	res, err := stmt.Exec("1", "xi", "C")
	checkErr(err)

	//pg不支持这个函数，因为他没有类似MySQL的自增ID
	// id, err := res.LastInsertId()
	// checkErr(err)
	// fmt.Println(id)

	var lastInsertId int
	err = db.QueryRow("INSERT INTO person(id,name,status)  VALUES($1,$2,$3) returning id;", "2", "li", "C").Scan(&lastInsertId)
	checkErr(err)
	fmt.Println("最后插入id =", lastInsertId)

	//更新数据
	stmt, err = db.Prepare("update person set name=$1 where id=$2")
	checkErr(err)

	res, err = stmt.Exec("zhijiang", 1)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)

	//查询数据
	rows, err := db.Query("SELECT id,name,status FROM person")
	checkErr(err)

	for rows.Next() {
		var id int
		var name string
		var status string
		err = rows.Scan(&id, &name, &status)
		checkErr(err)
		fmt.Println(id)
		fmt.Println(name)
		fmt.Println(status)
	}

	//删除数据
	stmt, err = db.Prepare("delete from person where id=$1")
	checkErr(err)

	res, err = stmt.Exec(1)
	checkErr(err)

	affect, err = res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)

	db.Close()

}

func checkDErr(err error) {
	if err != nil {
		panic(err)
	}
}
