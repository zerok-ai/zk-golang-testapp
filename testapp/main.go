package main

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/v12"
	"io"
	"net/http"
)

type AppContext struct {
	app *iris.Application
}

func (a AppContext) NewAppContext() *AppContext {
	return &AppContext{
		app: iris.New(),
	}
}

func main() {
	appCtx := AppContext{}
	app := appCtx.NewAppContext().app

	booksAPI := app.Party("/books")
	{
		booksAPI.Use(iris.Compression)

		// GET: http://localhost:8080/books
		booksAPI.Get("/", appCtx.list)
		booksAPI.Get("/from/{target}", appCtx.getOtherBooks)
	}

	app.Listen(":8080")
}

// Book example.
type Book struct {
	Title string `json:"title"`
}

func (a AppContext) getOtherBooks(ctx iris.Context) {
	println("getOtherBooks")
	var p struct {
		Target string `json:"target"`
	}
	ctx.ReadParams(&p)

	target := p.Target
	println(target)
	url := "http://" + target + "/books" // Replace with the actual URL.

	// Make the HTTP GET request.
	response, err := http.Get(url)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.WriteString(fmt.Sprintf("Failed to make the HTTP GET request: %s", err.Error()))
		return
	}
	defer response.Body.Close()

	// Check if the response status code is OK (200).
	if response.StatusCode != http.StatusOK {
		ctx.StatusCode(http.StatusBadGateway)
		ctx.WriteString(fmt.Sprintf("HTTP GET request failed with status code: %d", response.StatusCode))
		return
	}

	// Read the response body.
	body, err := io.ReadAll(response.Body)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.WriteString(fmt.Sprintf("Failed to read response body: %s", err.Error()))
		return
	}

	// Unmarshal the JSON payload into []Books.
	var books []Book
	err = json.Unmarshal(body, &books)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.WriteString(fmt.Sprintf("Failed to unmarshal JSON payload: %s", err.Error()))
		return
	}
	books = append(books, Book{"Go Go Go"})
	// Send the response body to the client.
	resp, err := json.Marshal(books)

	ctx.StatusCode(http.StatusOK)
	ctx.Write(resp)
}

func (a AppContext) list(ctx iris.Context) {
	println("list")
	books := []Book{
		{"Mastering Concurrency in Go"},
		{"Go Design Patterns"},
		{"Black Hat Go"},
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(books)
}
