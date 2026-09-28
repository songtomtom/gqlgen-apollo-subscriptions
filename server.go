package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	coderws "github.com/coder/websocket"
	"github.com/rs/cors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/songtomtom/gqlgen-apollo-subscriptions/graph"
	"github.com/songtomtom/gqlgen-apollo-subscriptions/graph/model"
)

const (
	defaultPort = "8080"
	// docker-compose.yml 의 MySQL 설정과 맞춘 기본값. 운영에서는 DSN 환경 변수로 덮어쓴다.
	defaultDSN = "test:test@tcp(127.0.0.1:33006)/test?charset=utf8mb4&parseTime=True&loc=Local"
)

func main() {
	port := envOr("PORT", defaultPort)
	dsn := envOr("DSN", defaultDSN)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	if err = db.AutoMigrate(&model.Post{}, &model.Comment{}); err != nil {
		log.Fatalf("failed to auto migrate schema: %v", err)
	}

	// 리졸버 의존성은 프로세스에 하나만 둔다.
	// Observer 가 요청마다 새로 만들어지면 뮤테이션과 구독이 서로 다른 Observer 를 보게 되어 이벤트가 전달되지 않는다.
	resolver := &graph.Resolver{
		DB:       db,
		Observer: graph.NewObserver(),
	}
	schema := graph.NewExecutableSchema(graph.Config{Resolvers: resolver})

	// 하나의 핸들러가 HTTP POST(Query, Mutation)와 WebSocket(Subscription)을 모두 받는다.
	// Apollo Client 는 graphql-ws 라이브러리로 graphql-transport-ws 서브프로토콜을 사용하고,
	// gqlgen 의 Websocket 전송은 이 프로토콜을 기본으로 지원한다.
	srv := handler.New(schema)
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.Websocket{
		// 브라우저의 WebSocket 은 CORS 를 타지 않으므로 오리진 검사를 여기서 직접 한다.
		// CRA 개발 서버(3000)와 playground(8080)만 허용한다. 운영에서는 실제 클라이언트 도메인으로 바꾼다.
		Implementation: transport.CoderWebsocketImplementation{
			AcceptOptions: coderws.AcceptOptions{
				OriginPatterns: []string{"localhost:3000", "localhost:" + port},
			},
		},
		// 프록시나 로드밸런서가 유휴 연결을 끊지 않도록 주기적으로 ping 을 보낸다.
		KeepAlivePingInterval: 10 * time.Second,
	})

	mux := http.NewServeMux()
	mux.Handle("/", playground.Handler("GraphQL playground", "/query"))
	// /query 와 /subscriptions 는 같은 핸들러다. 클라이언트가 HTTP 와 WS 주소를 나눠 쓰는 것을 그대로 지원한다.
	mux.Handle("/query", srv)
	mux.Handle("/subscriptions", srv)

	c := cors.New(cors.Options{AllowedOrigins: []string{"*"}})

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, c.Handler(mux)))
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
