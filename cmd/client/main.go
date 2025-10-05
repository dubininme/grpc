package main

import (
	"context"
	"log"

	"github.com/dubininme/grpc/pkg/api/example"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
		Content:  "Hello, world!",
		AuthorId: "author-123",
	})

	if err != nil {
		log.Fatal("could not create post: ", err)
	}

	log.Printf("Post created with ID: %d", resp.GetPostId())

	bytes, err := protojson.Marshal(resp)
	if err != nil {
		log.Fatal("could not marshal response: ", err)
	}

	log.Println("Response in JSON format:", string(bytes))
}
