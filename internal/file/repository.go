package file

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type Repository interface {
	Create(file File) error
	Update(file File) error
	Get(name string) (File, error)
	List(name string) ([]File, error)
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

func (r *fileRepositoryImpl) List(name string) ([]File, error) {
	var files []File
	query := `SELECT name, created_at, updated_at FROM files`
	err := r.db.Select(&files, query)
	return files, err
}
