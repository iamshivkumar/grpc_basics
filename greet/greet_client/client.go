package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/shivkumar123g/grpc_go_course/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Hello, I'm a client")
	certFile := "ssl/ca.crt"
	creds, sslErr := credentials.NewClientTLSFromFile(certFile, "")
	if sslErr != nil {
		log.Fatalf("Failed loading certificates: %v",sslErr)
		return
	}

	cc, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()
	c := greetpb.NewGreetServiceClient(cc)

	doUnary(c)
	// doServerStreaming(c)
	// doClientStreaming(c)
	// doBiDiStreaming(c)
	// doUnaryWithDeadline(c, 3*time.Second)
	// doUnaryWithDeadline(c, 2*time.Second)
	// doUnaryWithDeadline(c, 5*time.Second)
}

func doUnary(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Shivkumar",
			LastName:  "Konade",
		},
	}
	res, err := c.Greet(context.Background(), req)

	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)

	}
	log.Printf("Response from Greet: %v", res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Shivkumar",
			LastName:  "Konade",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GreetManyTimes RPC: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		log.Printf("Response from GreetManyTimes: %v", msg.GetResult())
	}

}

func doClientStreaming(c greetpb.GreetServiceClient) {

	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{FirstName: "Shivkumar", LastName: "Konade"},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{FirstName: "Rajkumar", LastName: "Konade"},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{FirstName: "Sachin", LastName: "Konade"},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{FirstName: "Ashwini", LastName: "Konade"},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{FirstName: "Tejeswini", LastName: "Konade"},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while calling LognGreet: %v", err)
	}
	for _, req := range requests {
		fmt.Printf("sent %v\n", req.GetGreeting().GetFirstName())
		stream.Send(req)
		time.Sleep(time.Second)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from LongGreet: %v", err)
	}

	fmt.Printf("LongGreet Response: %v\n", res)
}

func doBiDiStreaming(c greetpb.GreetServiceClient) {

	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{FirstName: "Shivkumar", LastName: "Konade"},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{FirstName: "Rajkumar", LastName: "Konade"},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{FirstName: "Sachin", LastName: "Konade"},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{FirstName: "Ashwini", LastName: "Konade"},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{FirstName: "Tejeswini", LastName: "Konade"},
		},
	}

	stream, err := c.GreetEveryone((context.Background()))
	if err != nil {
		log.Fatalf("error while calling GreetEveryone: %v", err)
	}
	waitc := make(chan struct{})
	go func() {
		for _, req := range requests {
			fmt.Printf("sent %v\n", req.GetGreeting().GetFirstName())
			stream.Send(req)
			time.Sleep(time.Second)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error while receiving: %v", err)
				break
			}
			fmt.Printf("response: %v", res.GetResult())
		}
		close(waitc)
	}()
	<-waitc
}

func doUnaryWithDeadline(c greetpb.GreetServiceClient, timeout time.Duration) {
	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Raj",
			LastName:  "Konade",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	res, err := c.GreetWithDeadline(ctx, req)

	if err != nil {

		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout was hit! Deadline was exceeded")
			} else {
				fmt.Printf("error: %v\n", err)
			}
		} else {
			log.Fatalf("error while calling Greet RPC: %v\n", err)
		}

		return
	}
	log.Printf("Response from Greet: %v\n", res.Result)
}
