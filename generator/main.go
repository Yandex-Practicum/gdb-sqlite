package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	sqlFile := `
CREATE TABLE IF NOT EXISTS clients (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	fio VARCHAR(128) NOT NULL DEFAULT "",
	login VARCHAR(32) NOT NULL DEFAULT "",
	birthday CHAR(8) NOT NULL DEFAULT "",
	email VARCHAR(64) NOT NULL DEFAULT ""
); 	

CREATE TABLE IF NOT EXISTS products (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	product VARCHAR(64) NOT NULL DEFAULT "",
	price INTEGER NOT NULL DEFAULT 0
);	

CREATE TABLE IF NOT EXISTS sales (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	client INTEGER NOT NULL DEFAULT 0,
	product INTEGER NOT NULL DEFAULT 0,
	volume INTEGER NOT NULL DEFAULT 1,
	date CHAR(8) NOT NULL DEFAULT ""
);	

INSERT INTO products (product, price) VALUES
	('Мои финансы', 750),
	('Суперпланировщик', 600),
	('Заметки', 450);

`
	clients, err := os.ReadFile("clients.csv")
	if err != nil {
		panic(err)
	}
	vals := make([]string, 0, 1000)
	for _, v := range strings.Split(string(clients), "\n") {
		fields := strings.Split(v, ",")
		if len(fields) != 6 {
			break
		}
		t, err := time.Parse("02.01.2006", fields[3])
		if err != nil {
			panic(err)
		}
		fio := strings.Join(fields[:3], " ")
		vals = append(vals, fmt.Sprintf("('%s','%s','%s','%s')", fio, t.Format("20060102"), fields[4], fields[5]))
	}
	sqlFile += "INSERT INTO clients (fio, login, birthday, email) VALUES\n"
	sqlFile += strings.Join(vals, ",\n") + ";\n"

	rnd := rand.New(rand.NewSource(7))
	vals = vals[:0]
	start := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 1; i <= 250; i++ {
		volume := rnd.Int31n(5) + 1
		if volume <= 3 {
			volume = 1
		} else {
			volume = rnd.Int31n(10) + 1
		}
		day := start.AddDate(0, 0, int(rnd.Int31n(320)+1))
		vals = append(vals, fmt.Sprintf("(%d,%d,%d,'%s')", i,
			rnd.Int31n(3)+1, volume, day.Format("20060102")))
		if rnd.Int31n(3) > 1 {
			i--
		}
	}
	sqlFile += "INSERT INTO sales (client, product, volume, date) VALUES\n"
	sqlFile += strings.Join(vals, ",\n") + ";\n"

	if err = os.WriteFile("demo.sql", []byte(sqlFile), 0644); err != nil {
		panic(err)
	}
}
