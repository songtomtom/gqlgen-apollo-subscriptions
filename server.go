package main

import (
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	"github.com/songtomtom/gqlgen-apollo-subscriptions/graph/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/songtomtom/gqlgen-apollo-subscriptions/graph"
)

const defaultPort = "8080"

var observer map[string]chan *model.Comment

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	dsn := "test:test@tcp(127.0.0.1:33006)/test?charset=utf8mb4&parseTime=True&loc=Local"

	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err = db.AutoMigrate(&model.Post{}, &model.Comment{}); err != nil {
		log.Fatalf("failed to auto migration schema: %v", err)
	}

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})

	observer = map[string]chan *model.Comment{}

	mux := http.NewServeMux()

	mux.Handle("/", playground.Handler("GraphQL playground", "/query"))
	mux.HandleFunc("/query", query(db, observer))
	mux.HandleFunc("/subscriptions", subscription(db, observer))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, c.Handler(mux)))
}

func query(db *gorm.DB, observer map[string]chan *model.Comment) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		srv := handler.NewDefaultServer(
			graph.NewExecutableSchema(
				graph.Config{
					Resolvers: &graph.Resolver{
						DB:       db,
						Observer: observer,
					},
				},
			),
		)
		srv.ServeHTTP(w, r)
	}
}

func subscription(db *gorm.DB, observer map[string]chan *model.Comment) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		srv := handler.New(
			graph.NewExecutableSchema(
				graph.Config{
					Resolvers: &graph.Resolver{
						DB:       db,
						Observer: observer,
					},
				},
			),
		)
		srv.AddTransport( // <---- This is the important part!
			&transport.Websocket{
				Upgrader: websocket.Upgrader{
					CheckOrigin: func(r *http.Request) bool {
						return true
					},
					ReadBufferSize:  1024,
					WriteBufferSize: 1024,
				},
			},
		)
		srv.ServeHTTP(w, r)
	}
}
