test_unit:
	go test ./... -v

test_enviroment:
	docker compose -f "scripts/e2e/docker-compose.yml" up -d

test_e2e: test_enviroment
	go test -tags=e2e ./test_e2e -v

