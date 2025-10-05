package main

import (
	"context"
	"io"
	"log"
	"math/rand/v2"
	"net"
	"sync"

	"buf.build/go/protovalidate"
	"github.com/dubininme/grpc/pkg/api/example"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

func main() {
	validator, err := protovalidate.New()
	if err != nil {
		log.Fatal("failed to create validator: ", err)
	}

	server := grpc.NewServer()
	service := &ExampleService{
		storage:  make(map[uint64]Post, 1),
		vaidator: validator,
	}

	example.RegisterExampleServer(server, service)

	lis, err := net.Listen("tcp", ":8085")
	if err != nil {
		log.Fatal("failed to listen: ", err)
	}

	reflection.Register(server)

	log.Println("gRPC server listening on :8085")
	if err := server.Serve(lis); err != io.EOF {
		log.Fatal("failed to serve: ", err)
	}
}

type Post struct {
	ID       uint64
	Title    string
	Content  string
	AuthorId string
}

type ExampleService struct {
	example.UnimplementedExampleServer

	vaidator *protovalidate.Validator
	storage  map[uint64]Post
	mx       sync.RWMutex
}

func (s *ExampleService) CreatePost(ctx context.Context, req *example.CreatePostRequest) (*example.CreatePostResponse, error) {
	if err := s.vaidator.Validate(req); err != nil {
		return nil, err
	}

	id := rand.Uint64()

	post := &Post{
		ID:       id,
		Title:    req.GetTitle(),
		Content:  req.GetContent(),
		AuthorId: req.GetAuthorId(),
	}

	s.mx.Lock()
	s.storage[id] = *post
	s.mx.Unlock()

	return &example.CreatePostResponse{
		PostId: id,
	}, nil
}

func (s *ExampleService) ListPosts(ctx context.Context, req *example.ListPostsRequest) (*example.ListPostsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListPosts not implemented")
}
