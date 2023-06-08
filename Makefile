PATHS = $(shell go list ./... | tail -n +2 | grep -v testintegration)

install_dependencies:
	go install gotest.tools/gotestsum@latest
	go install github.com/vektra/mockery/v2@v2.20.0
	go install github.com/ditrit/badaas/tools/badctl@latest

lint:
	golangci-lint run

test_unit_badaas:
	gotestsum --format pkgname $(PATHS)

test_unit_badctl:
	gotestsum --format pkgname ./tools/badctl/...

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

sqlserver:
	docker compose -f "docker/sqlserver/docker-compose.yml" up -d --build
	./docker/wait_for_db.sh

test_integration_postgresql: postgresql
	DB=postgresql gotestsum --format testname ./testintegration

test_integration_cockroachdb: cockroachdb
	DB=postgresql gotestsum --format testname ./testintegration

test_integration_mysql: mysql
	DB=mysql gotestsum --format testname ./testintegration -tags=mysql

test_integration_sqlite:
	DB=sqlite gotestsum --format testname ./testintegration

test_integration_sqlserver: sqlserver
	DB=sqlserver gotestsum --format testname ./testintegration

test_integration: test_integration_postgresql

test_e2e:
	docker compose -f "docker/cockroachdb/docker-compose.yml" -f "docker/test_api/docker-compose.yml" up -d
	./docker/wait_for_api.sh 8000/info
	go test ./test_e2e -v

test_generate_mocks:
	go generate

.PHONY: test_unit test_integration test_e2e

