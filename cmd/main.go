package main

import (
	"file_storage/internal/db"
	"file_storage/internal/handler"
	"file_storage/internal/limiter"
	pb "file_storage/pkg/grpc"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

const FilesListLimit = 100
const FilesUploadLimit = 10

func main() {
	err := db.Init()
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
	pb.RegisterFileServiceServer(s, &handler.ServerUpload{
		Limiter: limiter.NewLimiter(FilesUploadLimit),
	})

	pb.RegisterFileListServiceServer(s, &handler.ServerList{
		Limiter: limiter.NewLimiter(FilesListLimit),
	})

	log.Println("Server started at ", os.Getenv("GRPC_HOST"))
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
