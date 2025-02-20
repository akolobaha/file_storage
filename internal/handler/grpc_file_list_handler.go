package handler

import (
	"context"
	"file_storage/internal/db"
	f "file_storage/internal/file"
	"file_storage/internal/limiter"
	pb "file_storage/pkg/grpc"
)

type ServerList struct {
	pb.UnimplementedFileListServiceServer
	Limiter *limiter.Limiter
}

func (e ServerList) FilesList(ctx context.Context, empty *pb.Empty) (*pb.MultipleFile, error) {
	e.Limiter.Acquire()
	defer e.Limiter.Release()

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
