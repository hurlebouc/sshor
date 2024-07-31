FROM golang:1.22.5-alpine3.20 as base
run apk add --no-cache make

FROM base as dev

run apk add --no-cache git
run apk add --no-cache openssh

FROM base as build

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN make build-all

FROM scratch as package

COPY --from=build /usr/src/app/sshor-* /target/bin/
 CMD [ "executable" ]
