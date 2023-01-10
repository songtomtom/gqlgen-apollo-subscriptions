package main

import (
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/websocket"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/songtomtom/gqlgen-apollo-subscriptions/graph"
)

const defaultPort = "8080"

type (
	Post struct {
		gorm.Model
	}

	Comment struct {
		gorm.Model
		PostID  int
		Content string
	}
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := "test:test@tcp(127.0.0.1:33006)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err = db.AutoMigrate(&Post{}, &Comment{}); err != nil {
		log.Fatalf("failed to auto migration schema: %v", err)
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
