PATHS = $(shell go list ./... | tail -n +2 | grep -v testintegration)

lint:
	golangci-lint run

test_unit_badaas:
	gotestsum --format pkgname $(PATHS)

test_unit_badctl:
	go test ./tools/badctl/... -v

test_unit: test_unit_badaas test_unit_badctl

rmdb:
	docker stop badaas-test-db && docker rm badaas-test-db

test_db:
	docker compose -f "docker/test_db/docker-compose.yml" up -d
	./docker/wait_for_api.sh 8080/health

postgresql:
	docker compose -f "docker/postgresql/docker-compose.yml" up -d
	./docker/wait_for_db.sh

cockroachdb:
	docker compose -f "docker/cockroachdb/docker-compose.yml" up -d
	./docker/wait_for_db.sh

mysql:
	docker compose -f "docker/mysql/docker-compose.yml" up -d
	./docker/wait_for_db.sh

test_integration_postgresql: postgresql
	DB=postgresql go test ./testintegration -v

test_integration_cockroachdb: cockroachdb
	DB=postgresql go test ./testintegration -v

test_integration_mysql: mysql
	DB=mysql go test ./testintegration -v

test_integration: test_integration_postgresql

test_e2e:
	docker compose -f "docker/cockroachdb/docker-compose.yml" -f "docker/test_api/docker-compose.yml" up -d
	./docker/wait_for_api.sh 8000/info
	go test ./test_e2e -v

test_generate_mocks:
	go generate

.PHONY: test_unit test_integration test_e2e

