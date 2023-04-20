PATHS = $(shell go list ./... | tail -n +2 | grep -v test_integration | grep -v test_e2e)

test_unit:
	go test $(PATHS) -v

test_unit_and_cover:
	go test $(PATHS) -coverpkg=./... -coverprofile=coverage_unit.out -v

test_integration:
	docker compose -f "docker/test_db/docker-compose.yml" up -d
	go test ./test_integration -v

test_integration_and_cover:
	docker compose -f "docker/test_db/docker-compose.yml" up -d
	go test ./test_integration -coverpkg=./... -coverprofile=coverage_int.out -v

test_e2e:
	docker compose -f "docker/test_db/docker-compose.yml" -f "docker/test_api/docker-compose.yml" up -d
	./docker/wait_for_api.sh 8000/info
	go test ./test_e2e -v

example_birds:
	EXAMPLE=birds docker compose -f "docker/api/docker-compose.yml" up

example_posts:
	EXAMPLE=posts docker compose -f "docker/api/docker-compose.yml" up

badaas:
	docker compose -f "docker/api/docker-compose.yml" up

.PHONY: test_unit test_integration test_e2e example_birds example_posts badaas

