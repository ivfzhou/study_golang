package main

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
)

func main() {}

var db *sql.DB

func Connect() {
	cfg := mysql.Config{
		User:                 "ivfzhou",
		Passwd:               "123456",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "db_sql",
		AllowNativePasswords: true,
		Collation:            "utf8mb4_unicode_ci",
	}

	// db, err := sql.Open("mysql", cfg.FormatDSN())

	cons, _ := mysql.NewConnector(&cfg)
	db = sql.OpenDB(cons)

	_ = db.Ping()
}

type Album struct {
	ID     int
	Title  string
	Artist string
	Price  float64
}

func Create() {
	var alb Album
	result, _ := db.Exec("INSERT INTO `t_album` (`title`, `artist`, `price`) VALUES (?, ?, ?)", alb.Title, alb.Artist, alb.Price)
	id, _ := result.LastInsertId()
	_ = id
}

func Retrieve1() {
	var name string
	rows, _ := db.Query("SELECT * FROM `t_album` WHERE `artist` = ?", name)
	defer func() {
		_ = rows.Close()
	}()
	var albums []*Album
	for rows.Next() {
		var alb Album
		_ = rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price)
		albums = append(albums, &alb)
	}
}

func Retrieve2() {
	var id int
	row := db.QueryRow("SELECT * FROM `t_album` WHERE `id` = ?", id)
	var alb Album
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			return
		}
	}
}

func Retrieve3() {
	var s sql.NullString
	_ = db.QueryRow("SELECT `artist` FROM `t_album` WHERE `price` IS NULL").Scan(&s)

	if s.Valid {
		_ = s.String
	}
}

func Retrieve4() {
	rows, _ := db.Query("SELECT `price` FROM `t_album` WHERE `id` = 1; SELECT `price` FROM `t_album` WHERE `id` = 2;")
L1:
	for rows.Next() {
		a := &Album{}
		_ = rows.Scan(&a.Price)
		_ = rows.Err()
	}
	if rows.NextResultSet() {
		goto L1
	}
}

func Retrieve5() {
	stmt, _ := db.Prepare("SELECT `id`, `artist` from `t_album` WHERE `id` > ?;")
	var a Album
	_ = stmt.QueryRow(1).Scan(&a.ID, &a.Artist)

	tx, _ := db.Begin()
	stmt, _ = tx.Prepare("SELECT `artist` from `t_album` WHERE `id` > ?;")
	_ = stmt.QueryRow(4).Scan(&a.Artist)
	tx.Commit()
	tx.Rollback()
}
