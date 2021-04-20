package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/shivkumar123g/grpc_go_course/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	cc, err := grpc.Dial("0.0.0.0:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	c := calculatorpb.NewCalculatorServiceClient(cc)
	defer cc.Close()
	// doUnary(c)
	// doServerStreaming(c)
	// doClientStreaming(c)
	doErrorUnary(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	res, err := c.Sum(
		context.Background(), &calculatorpb.SumRequest{
			FistNumber:   5,
			SecondNumber: 45,
		},
	)
	if err != nil {
		log.Fatalf("Error while calling Sum RPC: %v", err)
	}
	fmt.Printf("Sum result is %v", res.GetSumResult())
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	resStream, err := c.PrimeNumberDecompsition(context.Background(), &calculatorpb.PrimeNumberDecompsitionRequest{
		Number: 12,
	})
	if err != nil {
		log.Fatalf("error while calling PrimeNumberDecompsition RPC")
	}
	for {
		res, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		fmt.Println(res.GetPrimeFactor())
	}
}

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	numbers := []int64{1, 2, 3, 4, 5, 6}
	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("error while calling PrimeNumberDecompsition RPC")
	}
	for _, v := range numbers {
		err := stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: v,
		})

		if err != nil {
			log.Fatalf("error while sending stream: %v", err)
		}
		time.Sleep(time.Second)

	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while reciving responce from server: %v", err)
	}
	fmt.Printf("average: %v", res.GetNumber())

}

func doErrorUnary(c calculatorpb.CalculatorServiceClient) {
	res, err := c.SquareRoot(
		context.Background(), &calculatorpb.SquareRootRequest{
			Number: -25,
		},
	)
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			fmt.Println(respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Printf("we probably sent a negative number\n")
			}
			return
		} else{
			log.Fatalf("Big Error calling SquareRoot: %v\n",err)
			return
		}
	} 
	fmt.Printf("Square root is %v\n", res.GetNumberRoot())
	
}
