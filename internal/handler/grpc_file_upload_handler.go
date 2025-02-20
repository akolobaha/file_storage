package handler

import (
	"file_storage/internal/db"
	f "file_storage/internal/file"
	"file_storage/internal/limiter"
	pb "file_storage/pkg/grpc"
	"io"
	"log"
	"os"
)

type ServerUpload struct {
	pb.UnimplementedFileServiceServer
	*limiter.Limiter
}

func (s *ServerUpload) UploadFile(stream pb.FileService_UploadFileServer) error {
	s.Limiter.Acquire()
	defer s.Limiter.Release()

	var file *os.File
	var filename string

	// Чтение данных из потока
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// Сохраним модель в бд
			repo := f.NewRepository(db.DB)
			service := f.NewService(repo)
			_, err := service.CreatOrUpdate(filename)

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
