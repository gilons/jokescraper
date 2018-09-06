package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

//Db is a global variable to hold the opened database connection
var Db *sql.DB

//Password global variable tohold the db password of the user
var Password string

func initialise() {

	db, err := sql.Open("mysql", "root:"+Password+"@tcp(127.0.0.1:3306)/JOKES")
	Db = db
	if err != nil {
		fmt.Println(err.Error())
	}
	sqlString := "CREATE TABLE `jokes_table` (" +
		"`joke_id` INT  PRIMARY KEY AUTO_INCREMENT NOT NULL, " +
		"`joke` TEXT ," +
		"`joke_info` VARCHAR(256), " +
		"`joke_title` VARCHAR(256) ," +
		"`category` VARCHAR(256) )"
	result, err := Db.Exec(sqlString)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(result)
}

func prepareDB() {
	msql := "CREATE DATABASE `JOKES`"
	db, err := sql.Open("mysql", "root:"+Password+"@tcp(127.0.0.1:3306)/")
	if err != nil {
		fmt.Println(err.Error())
	}
	result, err := db.Exec(msql)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(result)
}

func init() {
	arguments := os.Args
	Password = arguments[1]
	prepareDB()
	initialise()
}

//PushJokesToDB function pusshes the scraped content to the db
func PushJokesToDB(jokeTitle, joke, author, category, jokeInfo string) {
	sqlString := "INSERT INTO `jokes_table` SET " +
		"`joke_title` = \"" + jokeTitle + "\"" +
		", `joke_info` = \"" + jokeInfo + "\"" +
		", `joke` = \"" + joke + "\"" +
		", `category` = \"" + category + "\""
	result, err := Db.Exec(sqlString)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(result)
}
