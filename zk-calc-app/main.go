package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
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

func (s *calculatorServer) Add(ctx context.Context, req *calculator.AddRequest) (*calculator.AddResponse, error) {
	log.Println("found traceParent ", extractTraceParent(ctx))

	result := req.A + req.B
	fmt.Printf("Add result in destination: %d\n", result)
	return &calculator.AddResponse{Result: result}, nil
}
func (s *calculatorServer) CallOtherAdd(ctx context.Context, req *calculator.AddRequest) (*calculator.AddResponse, error) {
	forward := req.Forward
	traceparent := extractTraceParent(ctx)
	log.Println("found traceParent ", traceparent)
	conn, err := grpc.Dial(forward+".zk-calc-app.svc.cluster.local:50051", grpc.WithInsecure())
	//conn, err := grpc.Dial("localhost:"+forward, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := calculator.NewCalculatorClient(conn)

	// Call the Add operation
	md := metadata.New(map[string]string{
		"traceparent": traceparent, // Replace with your actual access token
	})

	outCtx := metadata.NewOutgoingContext(context.Background(), md)

	addResponse, err := client.Add(outCtx, req)
	if err != nil {
		log.Fatalf("Add request failed: %v", err)
	}
	fmt.Printf("Add result in source: %d\n", addResponse.Result)
	return addResponse, nil
}

func (s *calculatorServer) Subtract(ctx context.Context, req *calculator.SubtractRequest) (*calculator.SubtractResponse, error) {
	log.Println("found traceParent ", extractTraceParent(ctx))
	result := req.A - req.B
	fmt.Printf("Subtract result in destination: %d\n", result)
	return &calculator.SubtractResponse{Result: result}, nil
}

func (s *calculatorServer) CallOtherSubtract(ctx context.Context, req *calculator.SubtractRequest) (*calculator.SubtractResponse, error) {
	forward := req.Forward
	traceparent := extractTraceParent(ctx)
	log.Println("found traceParent ", traceparent)
	conn, err := grpc.Dial(forward+".zk-calc-app.svc.cluster.local:50051", grpc.WithInsecure())
	//conn, err := grpc.Dial("localhost:"+forward, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := calculator.NewCalculatorClient(conn)

	// Call the Add operation
	md := metadata.New(map[string]string{
		"traceparent": traceparent, // Replace with your actual access token
	})
	outCtx := metadata.NewOutgoingContext(context.Background(), md)
	subtractResponse, err := client.Subtract(outCtx, req)
	if err != nil {
		log.Fatalf("Subtract request failed: %v", err)
	}
	fmt.Printf("Subtract result in source: %d\n", subtractResponse.Result)
	return subtractResponse, nil
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
