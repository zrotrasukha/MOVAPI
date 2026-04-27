include .envrc

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help: 
	@echo "usage"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ":" | sed -e 's/^/ /'

.PHONY: confirm
confirm: 
	@echo -n 'Are you sure? [y|N] ' && read ans && [ $${ans:-N} = y ]

# ================================================================================== #
# DEVELOPMENT
# ================================================================================== #
## run/api: run the cmd/api application
.PHONY: run
run/api:
	@go run ./cmd/api -db-dsn ${DSN}

## db/psql: connect to the database using psql
.PHONY: db
db/psql:
	psql $$dsn

## db/migration/new name=$1: create a new ddatabase migration file with the given name
.PHONY: db/migration/new
db/migrations/new:
	@echo 'Creating migration files for $(name)...'
	migrate create -ext sql -dir ./migrations/ -seq ${name}

## db/migrations/up: apply all up database migrations
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./migrations/ -database $$dsn up


# ================================================================================== #
# BUILD
# ================================================================================== #
## build/api: build the cmd/api application
.PHONY: build
build/api: 
	@echo 'Building cmd/api...'
	go build -ldflags='-s' -o ./bin/api ./cmd/api
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux/amd64/api ./cmd/api


# ================================================================================== #
# QUALITY CONTROL
# ================================================================================== #

## tidy: format all .go files and tidy module dependencies
.PHONY: tidy
tidy: 
	@echo 'Formatting .go files...'
	go fmt ./...
	@echo 'Tidying module dependencies...'
	go mod tidy
	go mod verify
	go mod vendor


## audit: run quality control checks
.PHONY: audit
audit: tidy
	@echo 'Checking module dependencies...'
	go mod tidy -diff
	go mod verify
	@echo 'Vetting code...'
	go vet ./...
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

