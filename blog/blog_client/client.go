package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/shivkumar123g/grpc_go_course/blog/blogpb"
	"google.golang.org/grpc"
	// "google.golang.org/grpc/credentials"
)

func main() {
	// certFile := "ssl/ca.crt"
	// creds, sslErr := credentials.NewClientTLSFromFile(certFile, "")
	// if sslErr != nil {
	// 	log.Fatalf("Failed loading certificates: %v", sslErr)
	// 	return
	// }

	// cc, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(creds))
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()
	c := blogpb.NewBlogServiceClient(cc)
	listBlog(c)

}

func createBlog(c blogpb.BlogServiceClient) {
	res, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{
		Blog: &blogpb.Blog{
			AutherId: "1",
			Title:    "Title 5",
			Content:  "Content 5",
		},
	})
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	fmt.Println(res)
}

func readBlog(c blogpb.BlogServiceClient) {
	res, err := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{
		BlogId: "608236bcf26ded5be685f4a8",//608236bcf26ded5be685f4a8
	})
	if err != nil {
		fmt.Printf("ERROR: can't read: %v",err)
	}
	fmt.Println(res)
}

func updateBlog(c blogpb.BlogServiceClient) {
	res, err := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{
		Blog: &blogpb.Blog{
			Id: "608236bcf26ded5be685f4a8",
			AutherId: "2",
			// Title:    "Title 2",
			// Content:  "Content 2",
		},
	})
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	fmt.Println(res)
}


func deleteBlog(c blogpb.BlogServiceClient) {
	res, err := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{
		BlogId: "608236bcf26ded5be685f4a8",//608236bcf26ded5be685f4a8
	})
	if err != nil {
		fmt.Printf("ERROR: can't read: %v",err)
	}
	fmt.Println(res)
}


func listBlog(c blogpb.BlogServiceClient) {  
	stream, err := c.ListBlog(context.Background(),&blogpb.ListBlogRequest{})
	if err != nil {
		log.Fatalf("error while calling GreetManyTimes RPC: %v", err)
	}
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		log.Println(msg.GetBlog())
	}
}