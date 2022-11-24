.PHONY: run-service
run-service:
	@export GET_PORT=8080 && export SET_PORT=8081 && go run ./cmd/main.go

.PHONY: run-nats
run-nats:
	@docker run --rm -p 4222:4222 -ti nats:latest

.PHONY: run-master
run-master:
	@export GET_PORT=8080 && export SET_PORT=8081 && export REPLICAS_AVAILABILITY=false && export PRIORITY=1 && go run ./cmd/main.go

.PHONY: run-replica
run-replica:
	@export GET_PORT=8082 && export SET_PORT=8083 && export REPLICAS_AVAILABILITY=true && export REPLICA_PORT=8080 && export PRIORITY=2 && go run ./cmd/main.go
