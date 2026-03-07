package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	db, err := sql.Open("sqlite3", "store.db")

	if err != nil {
		log.Fatal("Ошибка при открытии БД: ", err)

	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Не удалось подключиться к БД: ", err)
	}

	log.Println("Успешное подключение к store.db!")

	createTableSQL := `
    CREATE TABLE IF NOT EXISTS products (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        model TEXT,
        company TEXT,
        price INTEGER
    );`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal("Ошибка создания таблицы: ", err)
	}

	log.Println("Таблица 'products' успешно создана или уже существует.")

}
