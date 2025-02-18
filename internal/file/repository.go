package file

import "github.com/jmoiron/sqlx"

type Repository interface {
	Save(file File) error
	Get(name string) (File, error)
	List(name string) ([]File, error)
}

type fileRepositoryImpl struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &fileRepositoryImpl{db: db}
}

func (r *fileRepositoryImpl) Save(file File) error {
	_, err := r.db.Exec("INSERT INTO files(name, data) VALUES ($1, $2)", file.Name)
	if err != nil {
		return err
	}
	return nil
}

func (r *fileRepositoryImpl) Get(name string) (File, error) {
	file := File{}
	query := `SELECT name, created_at, updated_at FROM files WHERE name = $1`
	err := r.db.Get(&file, query, name)
	return file, err
}

func (r *fileRepositoryImpl) List(name string) ([]File, error) {
	var files []File
	query := `SELECT name, created_at, updated_at FROM files`
	err := r.db.Select(&files, query)
	return files, err
}
