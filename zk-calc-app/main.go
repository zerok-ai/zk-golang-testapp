package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"zk-calc-app/calculator"
)

type calculatorServer struct{}

func extractTraceParent(ctx context.Context) string {
	md, _ := metadata.FromIncomingContext(ctx)
	// Access specific headers, e.g., the "Authorization" header.
	traceParent := md["traceparent"]
	if len(traceParent) > 0 {
		return traceParent[0]
	}

	return ""
}

func createCalculatorClient(serviceName string) (calculator.CalculatorClient, *grpc.ClientConn, error) {
	conn, err := grpc.Dial(serviceName+".zk-calc-app.svc.cluster.local:80", grpc.WithInsecure())
	//conn, err := grpc.Dial("localhost:"+serviceName, grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}
	client := calculator.NewCalculatorClient(conn)
	return client, conn, nil
}

func getContext(ctx context.Context, traceparent string) context.Context {
	//traceparent is 00-ddddddddf403a450d05c41c8c095634c-8d289ea473e34cdb-01
	//split it into four parts
	//traceparent := "00-ddddddddf403a450d05c41c8c095634c-8d289ea473e34cdb-01"

	//md := metadata.New(map[string]string{
	//	"traceparent": traceparent, // Replace with your actual access token
	//})
	//
	//outCtx := metadata.NewOutgoingContext(ctx, md)
	//return outCtx
	return ctx
}

func (s *calculatorServer) Add(ctx context.Context, req *calculator.AddRequest) (*calculator.AddResponse, error) {
	traceparent := extractTraceParent(ctx)
	log.Println("found traceParent ", traceparent)
	forward := req.Forward

	result := int32(0)
	log.Println("forward: ", forward)
	if forward == "" {
		result = req.A + req.B
		fmt.Printf("Add result in destination: %d\n", result)
	} else {
		calcClient, conn, err := createCalculatorClient(forward)
		if err != nil {
			return nil, err
		}
		defer conn.Close()
		withoutForward := &calculator.AddRequest{A: req.A, B: req.B, Forward: ""}
		addResponse, err := calcClient.Add(getContext(ctx, traceparent), withoutForward)
		if err != nil {
			log.Printf("Forwarded add request failed: %v", err)
			return nil, err
		}
		result = addResponse.Result
	}

	return &calculator.AddResponse{Result: result}, nil
}

func (s *calculatorServer) Divide(ctx context.Context, req *calculator.DivideRequest) (*calculator.DivideResponse, error) {
	traceparent := extractTraceParent(ctx)
	log.Println("found traceParent ", traceparent)
	forward := req.Forward

	result := int32(0)
	log.Println("forward: ", forward)
	if forward == "" {
		result = req.A / req.B
		fmt.Printf("Divide result in destination: %d\n", result)
	} else {
		calcClient, conn, err := createCalculatorClient(forward)
		if err != nil {
			return nil, err
		}
		defer conn.Close()
		withoutForward := &calculator.DivideRequest{A: req.A, B: req.B, Forward: ""}
		divideResponse, err := calcClient.Divide(getContext(ctx, traceparent), withoutForward)
		if err != nil {
			log.Printf("Forwarded divide request failed: %v", err)
			return nil, err
		}
		result = divideResponse.Result
	}

	return &calculator.DivideResponse{Result: result}, nil
}

func (s *calculatorServer) Restricted(ctx context.Context, req *calculator.RestrictedRequest) (*calculator.RestrictedResponse, error) {
	traceparent := extractTraceParent(ctx)
	log.Println("found traceParent ", traceparent)
	forward := req.Forward

	result := int32(0)
	log.Println("forward: ", forward)
	if forward == "" {
		fmt.Printf("Restricted result in destination")
		err := status.Error(codes.InvalidArgument, "Invalid argument provided")
		return nil, err
	} else {
		calcClient, conn, err := createCalculatorClient(forward)
		if err != nil {
			return nil, err
		}
		defer conn.Close()
		withoutForward := &calculator.RestrictedRequest{A: req.A, Forward: ""}
		restrictedResponse, err := calcClient.Restricted(getContext(ctx, traceparent), withoutForward)
		if err != nil {
			log.Printf("Forwarded restricted request failed: %v", err)
			return nil, err
		}
		result = restrictedResponse.Result
	}

	return &calculator.RestrictedResponse{Result: result}, nil
}

func (s *calculatorServer) Error(ctx context.Context, req *calculator.ErrorRequest) (*calculator.ErrorResponse, error) {
	traceparent := extractTraceParent(ctx)
	log.Println("found traceParent ", traceparent)
	forward := req.Forward

	log.Println("forward: ", forward)
	if forward == "" {
		fmt.Printf("Error result in destination")
		err := status.Error(codes.Code(req.Code), "Some custom error")
		return nil, err
	} else {
		calcClient, conn, err := createCalculatorClient(forward)
		if err != nil {
			return nil, err
		}
		defer conn.Close()
		withoutForward := &calculator.ErrorRequest{Code: req.Code, Forward: ""}
		_, err = calcClient.Error(getContext(ctx, traceparent), withoutForward)
		if err != nil {
			log.Printf("Forwarded error request failed: %v", err)
			return nil, err
		}
	}

	return &calculator.ErrorResponse{}, nil
}

func (s *calculatorServer) Subtract(ctx context.Context, req *calculator.SubtractRequest) (*calculator.SubtractResponse, error) {
	//log.Println("found traceParent ", extractTraceParent(ctx))
	//result := req.A - req.B
	//fmt.Printf("Subtract result in destination: %d\n", result)
	//return &calculator.SubtractResponse{Result: result}, nil

	traceparent := extractTraceParent(ctx)
	log.Println("found traceParent ", traceparent)
	forward := req.Forward

	result := int32(0)
	log.Println("forward: ", forward)
	if forward == "" {
		result = req.A - req.B
		fmt.Printf("Subtract result in destination: %d\n", result)
	} else {
		calcClient, conn, err := createCalculatorClient(forward)
		if err != nil {
			return nil, err
		}
		defer conn.Close()
		withoutForward := &calculator.SubtractRequest{A: req.A, B: req.B, Forward: ""}
		subtractResponse, err := calcClient.Subtract(getContext(ctx, traceparent), withoutForward)
		if err != nil {
			log.Printf("Forwarded subtract request failed: %v", err)
			return nil, err
		}
		result = subtractResponse.Result
	}

	return &calculator.SubtractResponse{Result: result}, nil
}

func (s *calculatorServer) mustEmbedUnimplementedCalculatorServer() {}

func main() {
	port := os.Args[1]
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()
	reflection.Register(server)
	calculator.RegisterCalculatorServer(server, &calculatorServer{})

	fmt.Println("Calculator gRPC server started on port " + port)
	if err := server.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
