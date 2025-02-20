package file

import (
	"database/sql"
	pb "file_storage/pkg/grpc"
	"fmt"
	"github.com/jmoiron/sqlx"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Repository interface {
	Create(file File) error
	Update(file File) error
	Get(name string) (File, error)
	List() ([]*pb.File, error)
}

type fileRepositoryImpl struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &fileRepositoryImpl{db: db}
}

func (r *fileRepositoryImpl) Create(file File) error {
	_, err := r.db.Exec(`
        INSERT INTO files(name) VALUES ($1)
        ON CONFLICT (name) DO UPDATE SET updated_at = $2`, file.Name, time.Now())

	if err != nil {
		return err
	}
	return nil
}

func (r *fileRepositoryImpl) Update(file File) error {
	_, err := r.db.Exec("UPDATE files SET updated_at = $1 WHERE name = $2", time.Now(), file.Name)
	if err != nil {
		return err
	}
	return nil
}

func (r *fileRepositoryImpl) Get(name string) (File, error) {
	file := File{}
	fmt.Println(r.db)
	query := `SELECT name, created_at, updated_at FROM files WHERE name = $1`
	err := r.db.Get(&file, query, name)
	if err != nil {
		return File{}, err
	}
	return file, nil
}

func (r *fileRepositoryImpl) List() ([]*pb.File, error) {
	var files []*pb.File

	query := `SELECT name, created_at, updated_at FROM files`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // Закрываем rows после завершения работы

	// Итерируем по строкам результата
	for rows.Next() {
		var file pb.File
		var createdAt, updatedAt sql.NullTime // Используем sql.NullTime для обработки возможных NULL значений

		// Сканируем значения в переменные
		if err := rows.Scan(&file.Name, &createdAt, &updatedAt); err != nil {
			return nil, err
		}

		// Преобразуем sql.NullTime в *timestamppb.Timestamp
		if createdAt.Valid {
			file.CreatedAt = timestamppb.New(createdAt.Time)
		}
		if updatedAt.Valid {
			file.UpdatedAt = timestamppb.New(updatedAt.Time)
		}

		// Добавляем файл в список
		files = append(files, &file)

	}

	return files, err
}
