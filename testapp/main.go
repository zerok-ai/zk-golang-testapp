package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	oteltrace "go.opentelemetry.io/otel/trace"
	"io/ioutil"
	"os"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"log"
	"net/http"
)

func printTraceFromContext(ctx context.Context) {
	trace := oteltrace.SpanContextFromContext(ctx)
	if trace.HasTraceID() {
		traceID := trace.TraceID()
		spanID := trace.SpanID()
		log.Printf("Before: TraceID: %s, SpanID: %s", traceID, spanID)
	} else {
		log.Printf("Before: TraceID: not found, SpanID: not found")
	}
}

func customMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// Log the trace information before the request is processed.
		printTraceFromContext(ctx)

		// Call the next handler in the chain.
		next.ServeHTTP(w, r)

		// Log the trace information after the request is processed.
		printTraceFromContext(ctx)
	})
}

func initTracer() (*sdktrace.TracerProvider, error) {
	// Create stdout exporter to be able to retrieve
	// the collected spans.

	//exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	//if err != nil {
	//	return nil, err
	//}

	collectorEndpoint := os.Getenv("ZEROK_EXPORTER_OTLP_ENDPOINT")
	if collectorEndpoint == "" {
		collectorEndpoint = "aws-otel-collector.aws-collector.svc.cluster.local"
	}
	exporter, err := otlptracehttp.New(context.Background(), otlptracehttp.WithEndpoint(collectorEndpoint+":4318"), otlptracehttp.WithInsecure())
	if err != nil {
		log.Fatalf("failed to create exporter: %v", err)
	}

	// For the demonstration, use sdktrace.AlwaysSample sampler to sample all traces.
	// In a production application, use sdktrace.ProbabilitySampler with a desired probability.
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceName("BookStore"))),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, err
}

func main() {
	tp, err := initTracer()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	http.Handle("/books", otelhttp.NewHandler(customMiddleware(http.HandlerFunc(list)), "list"))
	http.Handle("/books/from", otelhttp.NewHandler(customMiddleware(http.HandlerFunc(getOtherBooks)), "getOtherBooks"))

	// Start the server and listen on port 8080.
	fmt.Println("Starting server at port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

// Book example.
type Book struct {
	Title string `json:"title"`
}

func getOtherBooks(w http.ResponseWriter, r *http.Request) {
	//pritn header names 'traceparent'
	println(">> getOtherBooks::checking traceparent")
	printRequestHeaders(r, "getOtherBooks")

	println("getOtherBooks")
	target := r.URL.Query().Get("target")
	if target == "" {
		http.Error(w, "Target URL not specified", http.StatusBadRequest)
		return
	}
	url := "http://" + target + "/books" // Replace with the actual URL.
	println("target = " + target)
	println("url = " + url)

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

	var body []byte

	tr := otel.Tracer("Bookstore1")
	ctx := r.Context()
	bag, _ := baggage.Parse("username=donuts")
	reqCtx := baggage.ContextWithBaggage(ctx, bag)

	err := func(reqCtx context.Context) error {
		rCtx, span := tr.Start(reqCtx, "get books", oteltrace.WithAttributes(semconv.PeerServiceKey.String("BookStore1")))
		defer span.End()
		req, _ := http.NewRequestWithContext(rCtx, "GET", url, nil)

		fmt.Printf("Sending request...\n")
		res, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		body, err = ioutil.ReadAll(res.Body)
		_ = res.Body.Close()

		return err
	}(reqCtx)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
		return
	}

	// Unmarshal the JSON payload into []Books.
	var books []Book
	err = json.Unmarshal(body, &books)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		//w.WriteString(fmt.Sprintf("Failed to unmarshal JSON payload: %s", err.Error()))
		return
	}
	books = append(books, Book{"Go Go Go"})
	// Send the response body to the client.
	//resp, err := json.Marshal(books)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(books)
}

func list(w http.ResponseWriter, r *http.Request) {
	println(">> list::checking traceparent")
	printRequestHeaders(r, "list")

	println("list")
	books := []Book{
		{"Mastering Concurrency in Go"},
		{"Go Design Patterns"},
		{"Black Hat Go"},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(books)
}

func printRequestHeaders(r *http.Request, endpoint string) {
	//if r.Header != nil && r.Header["traceparent"] != nil {
	//	println("list::traceparent: " + r.Header["traceparent"][0])
	//} else {
	//	println("list::traceparent: not found")
	//}

	for name, headers := range r.Header {
		for _, value := range headers {
			fmt.Printf("%s::%s: %s\n", endpoint, name, value)
		}
	}
}
