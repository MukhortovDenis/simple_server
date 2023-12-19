FROM golang:latest AS BUILD
WORKDIR /build
COPY . .
RUN go build -o app ./cmd/main.go
ENTRYPOINT ["/build/app"]

#FROM alpine:latest AS DEPLOY
#WORKDIR /
#COPY --from=BUILD /build/app /app
#ENTRYPOINT ["/app"]
