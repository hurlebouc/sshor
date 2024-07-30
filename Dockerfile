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
RUN GOOS=linux go build -v -o /usr/local/bin/sshor-linux .
RUN GOOS=windows go build -v -o /usr/local/bin/sshor.exe .
RUN GOOS=darwin go build -v -o /usr/local/bin/sshor-darwin .

FROM scratch as package

COPY --from=build /usr/local/bin/sshor-linux /bin/sshor-linux
COPY --from=build /usr/local/bin/sshor.exe /bin/sshor.exe
COPY --from=build /usr/local/bin/sshor-darwin /bin/sshor-darwin

CMD ["/bin/sshor-linux"]
