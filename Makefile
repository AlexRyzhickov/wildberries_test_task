.PHONY: run-service
run-service:
	@export GET_PORT=8080 && export SET_PORT=8081 && go run ./cmd/main.go

.PHONY: run-nats
run-nats:
	@docker run --rm -p 4222:4222 -ti nats:latest