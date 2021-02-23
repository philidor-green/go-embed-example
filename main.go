package main

import (
	"fmt"
	"net/http"

	"embed"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

//go:embed embed/*

var content embed.FS

func main() {
	r := &Router{&mux.Router{}}

	r.MustResponse("GET", "/", func(res http.ResponseWriter, req *http.Request) {
		data, _ := content.ReadFile("embed/index.html")
		res.Header().Set("Content-Type", "text/html")
		fmt.Fprint(res, string(data))
	})

	r.MustResponse("GET", "/styles.css", func(res http.ResponseWriter, req *http.Request) {
		data, _ := content.ReadFile("embed/styles.css")
		res.Header().Set("Content-Type", "text/css")
		fmt.Fprint(res, string(data))
	})

	r.Run(":8080", "*")
}

type Router struct {
	*mux.Router
}

func (r *Router) MustResponse(meth, path string, h http.HandlerFunc) {
	r.HandleFunc(path, h).Methods(meth)
}

func (r *Router) Run(address, origins string) {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{origins},
		AllowedMethods:   []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "If-None-Match", "Content-Length", "Accept-Encoding", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)
	http.ListenAndServe(address, handler)
}

func vars(req *http.Request) map[string]string {
	return mux.Vars(req)
}
