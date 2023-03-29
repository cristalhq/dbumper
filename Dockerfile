FROM golang:1.17.3-alpine3.13 AS build_go

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -o dbumper.exec .

FROM alpine:3.17.3

COPY --from=build_go /app/dbumper.exec dbumper.exec
EXPOSE 8000

CMD ["./dbumper.exec"]
