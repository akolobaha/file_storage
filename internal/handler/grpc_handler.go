package handler

import (
	"context"
	"file_storage/internal/db"
	f "file_storage/internal/file"
	pb "file_storage/pkg/grpc"
	"fmt"
	"io"
	"log"
	"os"
)

type Server struct {
	pb.UnimplementedFileServiceServer
}

func (s *Server) UploadFile(stream pb.FileService_UploadFileServer) error {
	var file *os.File
	var filename string

	// Чтение данных из потока
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// Сохраним модель в бд
			repo := f.NewRepository(db.DB)
			service := f.NewService(repo)
			fileModel, err := service.CreatOrUpdate(filename)
			//err = service.SaveFile(fileModel)
			fmt.Println(fileModel)
			if err != nil {
				os.Remove("uploads/" + filename)
				return err
			}

			// Завершение потока
			return stream.SendAndClose(&pb.FileUploadResponse{
				Message: "File uploaded successfully",
				Status:  200,
			})
		}
		if err != nil {
			log.Printf("Failed to receive file chunk: %v", err)
			return err
		}

		// Если файл еще не создан, создаем новый файл на диске
		if file == nil {
			filename = req.GetFilename()
			file, err = os.Create("uploads/" + filename)
			if err != nil {
				log.Printf("Failed to create file: %v", err)
				return err
			}
		}

		// Записываем данные из чанка в файл
		_, err = file.Write(req.GetChunk())
		if err != nil {
			log.Printf("Failed to write chunk to file: %v", err)
			return err
		}
	}
}

func (e Server) FilesList(ctx context.Context, empty *pb.Empty) (*pb.MultipleFile, error) {

	repo := f.NewRepository(db.DB)
	service := f.NewService(repo)

	list, err := service.List()

	multipleFile := &pb.MultipleFile{
		Files: list,
	}

	if err != nil {
		return nil, err
	}

	return multipleFile, nil
}
