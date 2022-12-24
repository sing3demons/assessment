FROM golang:1.19-alpine as build

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go test --tags=unit -v ./...

RUN go build -o ./out/go-app .

FROM alpine:3.16.2
COPY --from=build /app/out/go-app /app/go-app

RUN adduser -u 1001 -D -s /bin/sh -g ping 1001
RUN chown 1001:1001 /app

RUN chmod +x /app
USER 1001

CMD ["/app/go-app"]