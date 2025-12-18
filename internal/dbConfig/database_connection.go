package dbConfig

import (
	"PilaiteProject/internal/db"
	"PilaiteProject/internal/wrapper"
	"context"
	_ "database/sql"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Connection struct {
	Queries *wrapper.AppQueries
	Pool    *pgxpool.Pool
}

func loadEnvFile(filename string) error {
	err := godotenv.Load(filename)
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	return err
}

func getDbURL() (databaseUrl string) {
	var dbHost = os.Getenv("DB_HOST")
	//fmt.Println("DB Host: " + dbHost)
	var dbPort = os.Getenv("DB_PORT")
	//fmt.Println("DB Port: " + dbPort)
	var dbUser = os.Getenv("DB_USER")
	//fmt.Println("User= " + dbUser)
	var dbPass = os.Getenv("DB_PASS")
	//fmt.Println("password= " + dbPass)
	var dbName = os.Getenv("DB_NAME")
	//fmt.Println("DB Name= " + dbName)

	databaseUrl = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName,
	)

	return databaseUrl
}

func ConnectDb(ctx context.Context) (*Connection, error) {
	err1 := loadEnvFile(".env")
	if err1 != nil {
		log.Println("Failed to load .env")
	}

	dbConn, err2 := pgxpool.New(ctx, getDbURL())
	if err2 != nil {
		log.Fatal(err2)
	}

	err3 := dbConn.Ping(ctx)
	if err3 != nil {
		log.Fatal("Cannot ping db", err3)
	}
	fmt.Println("Connected to database successfully")

	connection := &Connection{
		Queries: wrapper.NewAppQueries(db.New(dbConn)),
		Pool:    dbConn,
	}
	return connection, nil
}
