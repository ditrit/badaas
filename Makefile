test_unit:
	go test ./... -v

test_integration:
	docker compose -f "docker/test_db/docker-compose.yml" up -d
	go test -tags=integration ./test_integration -v

test_e2e:
	docker compose -f "docker/test_db/docker-compose.yml" -f "docker/test_api/docker-compose.yml" up -d
	./docker/wait_for_api.sh 8000/info
	go test -tags=e2e ./test_e2e -v

.PHONY: test_unit test_integration test_e2e

