NAME=fizzbuzz
BIN="bin/${NAME}"
MAIN=cmd/main.go

test:
	go test -failfast ./... -v -p=1 -count=1 -coverprofile .coverage.txt
	go tool cover -func .coverage.txt

build:
	go build -o ${BIN} ${MAIN}

run: build
	./${BIN}
