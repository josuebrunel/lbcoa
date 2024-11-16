# Build
FROM    golang:alpine AS build

ENV     GO111MODULE=on
ENV     CGO_ENABLED=1
ENV     GOOS=linux

RUN     apk update && apk upgrade
RUN     apk add --no-cache gcc musl-dev sqlite make
RUN     mkdir /go/src/app
WORKDIR /go/src/app
ADD     . /go/src/app
RUN     make deps
RUN     go mod tidy
RUN     make build
EXPOSE  8090


# Deploy
FROM    alpine
RUN     mkdir /opt/fizzbuzz
WORKDIR /opt/fizzbuzz
COPY    --from=build /go/src/app/bin/fizzbuzz /opt/fizzbuzz/fizzbuzz
EXPOSE  8090
CMD     ["./fizzbuzz"]
