docker-build:
	docker-compose up --build -d

run-leader-election:
	docker exec go-boojum_demo_1 go run ./cmd/demo
