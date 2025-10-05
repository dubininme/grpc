package main

import (
	"context"
	"log"

	"github.com/dubininme/grpc/pkg/api/example"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	conn, err := grpc.NewClient(
		":8085",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Fatal(err)
	}

	client := example.NewExampleClient(conn)

	resp, err := client.CreatePost(context.Background(), &example.CreatePostRequest{
		Title:    "My first post",
		AuthorId: "author-123",
		Content:  "Hello, world!",
	})

	if err != nil {
		switch status.Code(err) {
		case codes.InvalidArgument:
			log.Println("Invalid argument:")
		default:
			log.Fatal("CreatePost failed:", err)
		}

		if st, ok := status.FromError(err); ok {
			log.Println("code", st.Code(), "message", st.Message(), "details", st.Details())
		} else {
			log.Println("non gRPC error:", err)
		}
		return
	}

	log.Printf("Post created with ID: %d", resp.GetPostId())

	bytes, err := protojson.Marshal(resp)
	if err != nil {
		log.Fatal("could not marshal response: ", err)
	}

	log.Println("Response in JSON format:", string(bytes))
}
