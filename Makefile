NAME=fizzbuzz
BIN="bin/${NAME}"
MAIN=cmd/main.go
DBFILE=fizzbuzz.db
SWAG_GFILE=app/app.go
SWAG_OUT=app/docs

install-deps:
	go install github.com/swaggo/swag/cmd/swag@latest

swag.gen:
	swag init --parseInternal --parseDependency -g ${SWAG_GFILE} -output ${SWAG_OUT}

migrate:
	sqlite3 ${DBFILE} < migrations/0001_init.sql

test:
	go test -failfast ./... -v -p=1 -count=1 -coverprofile .coverage.txt
	go tool cover -func .coverage.txt

build: swag.gen
	go build -o ${BIN} ${MAIN}

run: build
	./${BIN}
