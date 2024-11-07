run:
	go run ./cmd/auth

run-with-jq:
	go run ./cmd/auth | jq

migrations-run:
	go run ./cmd/migrations

migrations-run-with-jq:
	go run ./cmd/migrations | jq

bench:
	go test -bench=. -count=5 ./tests

test:
	go test -v ./...