PATHS = $(shell go list ./... | tail -n +2)

lint:
	golangci-lint run

test_unit_badaas:
	go test $(PATHS) -v

test_unit_badaas_cover:
	go test $(PATHS) -coverpkg=./... -coverprofile=coverage_unit.out -v

test_unit_badctl:
	go test ./tools/badctl/... -v

test_unit: test_unit_badaas test_unit_badctl

test_db:
	docker compose -f "docker/test_db/docker-compose.yml" up -d

test_e2e:
	docker compose -f "docker/test_db/docker-compose.yml" -f "docker/test_api/docker-compose.yml" up -d
	./docker/wait_for_api.sh 8000/info
	go test ./test_e2e -v

test_generate_mocks:
	mockery --all --keeptree

.PHONY: test_unit test_e2e

