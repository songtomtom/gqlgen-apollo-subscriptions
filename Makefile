gen:
	go run -mod=mod github.com/99designs/gqlgen generate
start:
	go run server.go