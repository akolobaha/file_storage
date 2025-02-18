package main

import (
	"file_storage/internal/db"
	"file_storage/internal/handler"
	pb "file_storage/pkg/grpc" // Импортируйте сгенерированные файлы
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}

	err = db.Init()
	if err != nil {
		panic("failed to connect database")
	}
	defer func(DB *sqlx.DB) {
		DB.Close()
	}(db.DB)

	lis, err := net.Listen("tcp", os.Getenv("GRPC_HOST"))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterFileServiceServer(s, &handler.Server{})

	log.Println("Server started at ", os.Getenv("GRPC_HOST"))
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
