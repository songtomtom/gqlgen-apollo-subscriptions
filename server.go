package main

import (
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/songtomtom/gqlgen-apollo-subscriptions/graph"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	_, err := open("test", "test", "127.0.0.1", "33006", "test")
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.HandleFunc("/query", queryHandlerFn)
	http.HandleFunc("/subscriptions", subscriptionHandlerFn)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func queryHandlerFn(w http.ResponseWriter, r *http.Request) {
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	srv.ServeHTTP(w, r)
}

func subscriptionHandlerFn(w http.ResponseWriter, r *http.Request) {
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	srv.AddTransport(
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
	srv.ServeHTTP(w, r)
}
