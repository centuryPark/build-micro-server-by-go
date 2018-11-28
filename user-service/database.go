package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"fmt"
	"os"
)

func CreateConnection() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	fmt.Println(host)
	fmt.Println(user)
	fmt.Println(password)
	fmt.Println(dbName)

	return gorm.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s user=%s dbname=%s sslmode=disable password=%s",
			host, user, dbName, password,
		),
	)
}
