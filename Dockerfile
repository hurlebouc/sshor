FROM golang:1.22.5-alpine3.20 as base

FROM base as dev

run apk add --no-cache git
run apk add --no-cache openssh

FROM base as build

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN GOOS=linux GOARCH=amd64 go build -v -o /usr/local/bin/sshor-linux-amd64 .
RUN GOOS=linux GOARCH=386 go build -v -o /usr/local/bin/sshor-linux-386 .
RUN GOOS=windows GOARCH=amd64 go build -v -o /usr/local/bin/sshor-windows-amd64.exe .
RUN GOOS=windows GOARCH=386 go build -v -o /usr/local/bin/sshor-windows-386.exe .
RUN GOOS=darwin GOARCH=amd64 go build -v -o /usr/local/bin/sshor-darwin-amd64 .

FROM scratch as package

COPY --from=build /usr/local/bin /target/bin
 CMD [ "executable" ]
