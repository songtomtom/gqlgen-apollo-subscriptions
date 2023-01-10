gen:
	go run -mod=mod github.com/99designs/gqlgen generate
start:
	go run server.go
client:
	cd client && yarn start

up:
	docker-compose -f ./docker-compose.yml up --build -d
down:
	docker-compose -f ./docker-compose.yml down