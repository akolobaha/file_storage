package grpc_test

import (
	"context"
	"testing"
	"time"

	"file_storage/pkg/grpc"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

const serverAddr = "localhost:50051"

func TestUploadFile(t *testing.T) {
	// Set up connection to the server
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		t.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	client := grpc_gen.NewFileServiceClient(conn)

	// Simulate uploading a file in chunks
	fileChunks := []grpc_gen.FileChunk{
		{Chunk: []byte("file chunk 1"), Filename: "test_file.txt"},
		{Chunk: []byte("file chunk 2"), Filename: "test_file.txt"},
	}

	// Create a stream to send the chunks
	stream, err := client.UploadFile(context.Background())
	if err != nil {
		t.Fatalf("Error creating stream: %v", err)
	}

	// Send file chunks to the server
	for _, chunk := range fileChunks {
		if err := stream.Send(&chunk); err != nil {
			t.Fatalf("Failed to send chunk: %v", err)
		}
	}

	// Close the stream and get the response
	response, err := stream.CloseAndRecv()
	if err != nil {
		t.Fatalf("Failed to receive response: %v", err)
	}

	assert.Equal(t, "File uploaded successfully", response.Message)
	assert.Equal(t, int32(200), response.Status)
}

func TestFilesList(t *testing.T) {
	// Set up connection to the server
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		t.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	client := grpc_gen.NewFileListServiceClient(conn)

	// Call the FilesList endpoint
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	resp, err := client.FilesList(ctx, &grpc_gen.Empty{})
	if err != nil {
		t.Fatalf("Failed to fetch file list: %v", err)
	}

	// Validate that the response contains files
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Files)
}
