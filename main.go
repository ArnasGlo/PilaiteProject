package main

import (
	"PilaiteProject/internal/dbConfig"
	_ "PilaiteProject/internal/handler"
	"PilaiteProject/internal/server"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
)

func main() {
	ctx := context.Background()

	conn, err := dbConfig.ConnectDb(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Pool.Close()

	serverConfig := server.ServerConfig{
		Host: "localhost",
		Port: "8080",
	}

	server := server.NewServer(serverConfig, conn)
	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")

}
