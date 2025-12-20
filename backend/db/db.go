package db

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect() *sql.DB{
	data_url := os.Getenv("DATABASE_URL")
	if data_url == ""{
		log.Fatal("DATABASE_URL is not set..")
	}
	db,err := sql.Open("pgx",data_url)
	if err != nil{
		log.Fatalf("Failed to open database: %v",err)
	}
	/*---- The time limits -----*/
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(30*time.Minute)

	if err := db.Ping();
  err!=nil{
		log.Fatalf("Failed to ping to the server ... : %v",err)
	}
	log.Println("databse connection is successful ... ")
	return db
}