PATHS = $(shell go list ./... | tail -n +2 | grep -v testintegration | grep -v test_e2e)

test_unit:
	go test $(PATHS) -v

test_unit_and_cover:
	go test $(PATHS) -coverpkg=./... -coverprofile=coverage_unit.out -v

rmdb:
	docker stop badaas-test-db && docker rm badaas-test-db

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

