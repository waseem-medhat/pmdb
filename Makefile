include .env

gup:
	cd sql/schema && goose postgres ${DBURL} up

gdown:
	cd sql/schema && goose postgres ${DBURL} down

gstatus:
	cd sql/schema && goose postgres ${DBURL} status

sqlc:
	sqlc generate

run:
	go run ./cmd/pmdb

tw:
	pnpm run tw-watch

lint:
	golangci-lint run
