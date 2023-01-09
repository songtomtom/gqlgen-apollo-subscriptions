package main

import (
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/songtomtom/gqlgen-apollo-subscriptions/graph"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)
	http.HandleFunc("/subscriptions", func(w http.ResponseWriter, r *http.Request) {
		_srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
		_srv.AddTransport(
			&transport.Websocket{
				Upgrader: websocket.Upgrader{
					CheckOrigin: func(r *http.Request) bool {
						return true
					},
					ReadBufferSize:  1024,
					WriteBufferSize: 1024,
				},
			},
		) // <---- This is the important part!
		_srv.ServeHTTP(w, r)
	})

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
