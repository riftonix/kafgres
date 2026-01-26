.PHONY: build run test clean publish

build:
	go build -o kafgres ./cmd/kafgres

publish:
	KO_DOCKER_REPO=ghcr.io/riftonix ko build ./cmd/kafgres --base-import-paths --tags latest,v0.0.6

run: build
	./kafgres

test:
	go test ./...

clean:
	rm -f kafgres

up:
	docker-compose -f deploy/docker-compose.yml up -d

down:
	docker compose -f deploy/docker-compose.yml down -v

check-kafka:
	docker compose -f deploy/docker-compose.yml exec kafka kafka-console-consumer --bootstrap-server kafka:29092 --topic test-topic --from-beginning

check-kafgres:
	docker logs deploy-kafgres-1 && curl localhost:8080/health
