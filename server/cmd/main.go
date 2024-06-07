package main

import (
	"User/server/db"
	"log"
	"net"
	"net/rpc"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	database, err := db.ConnectDB()
	if err != nil {
		log.Println(err)
		log.Fatal("DB connection could not be set")
	}
	storage := db.NewStorage(database)
	userService := db.NewUserHandler(storage)
	rpc.Register(userService)

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("RPC server port 1234 da ishlamoqda")
	rpc.Accept(listener)
}
