package file

import (
	"database/sql"
	pb "file_storage/pkg/grpc"
	"github.com/go-faster/errors"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s Service) CreatOrUpdate(filename string) (File, error) {
	file, err := s.repo.Get(filename)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Файл не найден, создаем новый
			newFile := File{
				Name:      filename,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			// Возможно, нужно сохранить новый файл в базе данных
			if saveErr := s.repo.Create(newFile); saveErr != nil {
				return File{}, errors.Wrap(saveErr, "failed to save new file")
			}
			return newFile, nil
		}
		return File{}, errors.Wrap(err, "get file")
	}

	// Файл найден, обновляем время
	file.UpdatedAt = time.Now()
	err = s.repo.Update(file)
	if err != nil {
		return File{}, err
	}
	return file, nil
}

func (s Service) List() ([]*pb.File, error) {
	return s.repo.List()
}
