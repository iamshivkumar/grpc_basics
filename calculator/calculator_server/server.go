package main

import (
	"context"
	"io"
	"log"
	"math"
	"net"

	"github.com/shivkumar123g/grpc_go_course/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	a := req.GetFistNumber()
	b := req.GetSecondNumber()
	c := a + b
	return &calculatorpb.SumResponse{
		SumResult: c,
	}, nil
}

func (*server) PrimeNumberDecompsition(req *calculatorpb.PrimeNumberDecompsitionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompsitionServer) error {
	n := req.GetNumber()
	k := int64(2)
	for n > 1 {
		if n%k == 0 {
			stream.Send(&calculatorpb.PrimeNumberDecompsitionResponse{
				PrimeFactor: k,
			})
			n = n / k

		} else {
			k++
		}
	}
	return nil
}

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	v := int64(0)
	n := int64(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(
				&calculatorpb.ComputeAverageResponse{
					Number: float64(v) / float64(n),
				},
			)
		}
		if err != nil {
			log.Fatalf("error while reciving sream from client")
		}
		v += req.GetNumber()
		n++
	}
}

func (*server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Received a negative number: %v", number)
	}
	return &calculatorpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	calculatorpb.RegisterCalculatorServiceServer(s, &server{})
	
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
