package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"
    
	"github.com/shivkumar123g/grpc_go_course/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet fucntion was invoked with %v", req)
	fistName := req.GetGreeting().GetFirstName()

	result := "Hello " + fistName
	res := &greetpb.GreetResponse{
		Result: result,
	}
	return res, nil
}


func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		res := &greetpb.GreetManyTimesResponse{
			Result: "Hello " + firstName + " number " + strconv.Itoa(i),
		}
		stream.Send(res)
		time.Sleep(time.Second)
	}
	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	result := "Hello "
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client streamin: %v", err)
		}
		firstName := req.GetGreeting().GetFirstName()
		result += firstName + "! "
	}
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("error while reciving stream from client: %v", err)
			return err
		}
		sErr := stream.Send(&greetpb.GreetEveryoneResponse{
			Result: "Hello " + req.GetGreeting().GetFirstName() + "!\n",
		})
		if sErr != nil {
			log.Fatalf("Error while sending data to client: %v",sErr)
			return sErr
		}
	}
}

func (*server) GreetWithDeadline(ctx context.Context, req *greetpb.GreetWithDeadlineRequest) (*greetpb.GreetWithDeadlineResponse, error) {
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			return nil , status.Error(codes.Canceled,"The Client cancelled the request")
		}
		time.Sleep(time.Second)
	} 
	fmt.Printf("Greet fucntion was invoked with %v", req)
	fistName := req.GetGreeting().GetFirstName()

	result := "Hello " + fistName
	res := &greetpb.GreetWithDeadlineResponse{
		Result: result,
	}
	return res, nil
}


func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	certFile := "ssl/server.crt"
	keyFile := "ssl/server.pem"
	creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
	 
	if sslErr != nil {
		log.Fatalf("Failed loading certificates: %v",sslErr)
		return
	}


	s := grpc.NewServer(grpc.Creds(creds))
	greetpb.RegisterGreetServiceServer(s, &server{})
    
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
