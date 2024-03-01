go-run:
	go run cmd/server/main.go

docker-build:
	docker build -t victorhilario/lab-temperatura-cep:latest .

docker-run:
	docker run --rm -p 8080:8080 victorhilario/lab-temperatura-cep:latest

docker-compose:
	docker-compose up -d