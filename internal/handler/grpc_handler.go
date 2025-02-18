package handler

import (
	"context"
	pb "file_storage/pkg/grpc" // Импортируйте сгенерированные файлы
	"fmt"
	"io/ioutil"
)

type Server struct {
	pb.UnimplementedFileServiceServer
}

func (s *Server) UploadFile(ctx context.Context, req *pb.FileRequest) (*pb.FileResponse, error) {
	// Сохраните файл на диск
	fmt.Println("Upload File")
	err := ioutil.WriteFile(req.Filename, req.Bytes, 0644)
	if err != nil {
		return &pb.FileResponse{Message: "Failed to save file", Status: 1}, err
	}
	return &pb.FileResponse{Message: "File uploaded successfully", Status: 0}, nil

}
